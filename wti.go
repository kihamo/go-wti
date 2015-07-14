package godic

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	wti "github.com/fromYukki/webtranslateit_go_client"
)

const (
	updatedLangRegexp = "(?is)\\$lang\\[\"(?P<key>.*?)\"\\] = \"(?P<value>.*?)\";"
)

var (
	exp           *regexp.Regexp
	localeAliases = map[string]string{
		"en_en_vn": "en_vn",
		"en_en_th": "en_th",
		"en_en_id": "en_id",
		"ms":       "ms_my",
	}
)

func init() {
	var err error

	exp, err = regexp.Compile(updatedLangRegexp)
	if err != nil {
		panic(err)
	}
}

type WebTranslateIt struct {
	wti                 *wti.WebTranslateIt
	updateRetryDelay    time.Duration
	updateRetryAttempts int64
	dictionaries        []*Dictionary
	callback            func([]string)
}

type Dictionary struct {
	Locale  string
	Phrases map[string]string
	Hash    string
	Update  bool
}

func NewWebTranslateIt(wtiToken string, updateRetryDelay time.Duration, updateRetryAttempts int64) *WebTranslateIt {
	return &WebTranslateIt{
		wti:                 wti.NewWebTranslateIt(wtiToken),
		updateRetryDelay:    updateRetryDelay,
		updateRetryAttempts: updateRetryAttempts,
		dictionaries:        []*Dictionary{},
	}
}

func (w *WebTranslateIt) GetDictionary(locale string) (*Dictionary, error) {
	locale = w.getLocale(locale)

	for i := range w.dictionaries {
		if w.dictionaries[i].Locale == locale {
			return w.dictionaries[i], nil
		}
	}

	return nil, errors.New("Locale not exists")
}

func (w *WebTranslateIt) SetCallback(f func([]string)) {
	w.callback = f
}

func (w *WebTranslateIt) Update() error {
	for i := range w.dictionaries {
		w.dictionaries[i].Update = false
	}

	project, err := w.wti.GetProject()
	if err != nil {
		return err
	}

	zipFile, err := project.ZipFile()
	if err != nil {
		return err
	}

	data, err := zipFile.Extract()
	if err != nil {
		return err
	}

	for i := range project.ProjectFiles {
		w.parseFile(project.ProjectFiles[i], data[project.ProjectFiles[i].Name])
	}

	if w.callback != nil {
		locales := []string{}

		for i := range w.dictionaries {
			if w.dictionaries[i].Update {
				locales = append(locales, w.dictionaries[i].Locale)
			}
		}

		if len(locales) > 0 {
			w.callback(locales)
		}
	}

	if w.updateRetryDelay > 0 {
		time.AfterFunc(w.updateRetryDelay, func() {
			err := w.Update()
			if err != nil {
				log.Printf("Update error %s\n", err)
			}
		})
	}

	return nil
}

func (w *WebTranslateIt) getLocale(locale string) string {
	locale = strings.ToLower(locale)

	if alias, ok := localeAliases[locale]; ok {
		return alias
	}

	return locale
}

func (w *WebTranslateIt) parseFile(file wti.File, content []byte) {
	dictionary, err := w.GetDictionary(file.LocaleCode)
	if err != nil {
		dictionary = &Dictionary{
			Locale: w.getLocale(file.LocaleCode),
			Hash:   file.Hash,
		}
	}

	if dictionary.Hash != file.Hash {
		dictionary.Phrases = map[string]string{}
		dictionary.Update = true

		for _, match := range exp.FindAllStringSubmatch(string(content), -1) {
			value := match[2]
			if value == "" {
				value = match[1]
			}

			dictionary.Phrases[match[1]] = value
		}
	} else {
		dictionary.Update = false
	}
}
