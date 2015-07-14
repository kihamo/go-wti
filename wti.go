package godic

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	wti "github.com/fromYukki/webtranslateit_go_client"
)

const (
	updatedLangRegexp   = "(?is)\\$lang\\[\"(?P<key>.*?)\"\\] = \"(?P<value>.*?)\";"
	updateRetryDelay    = time.Second * 10
	updateRetryAttempts = 3
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
	wti          *wti.WebTranslateIt
	dictionaries []*Dictionary
	callback     func([]string)
	attempts     int64
	mutex        sync.Mutex
	updateAt     string
}

type Dictionary struct {
	Locale  string
	Phrases map[string]string
	Hash    string
	Update  bool
	File    wti.File
}

func GetLocale(locale string) string {
	locale = strings.ToLower(locale)

	if alias, ok := localeAliases[locale]; ok {
		return alias
	}

	return locale
}

func NewWebTranslateIt(wtiToken string) *WebTranslateIt {
	return &WebTranslateIt{
		wti:          wti.NewWebTranslateIt(wtiToken),
		dictionaries: []*Dictionary{},
	}
}

func (w *WebTranslateIt) GetDictionaries() []*Dictionary {
	return w.dictionaries
}

func (w *WebTranslateIt) GetDictionary(locale string) (*Dictionary, error) {
	locale = GetLocale(locale)

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

func (w *WebTranslateIt) Update() (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.attempts = 0
	for w.attempts < updateRetryAttempts {
		err = w.update()
		if err == nil {
			break
		} else if w.attempts <= updateRetryAttempts {
			time.Sleep(updateRetryDelay)
		}

		log.Printf("Update failed. Attempts: %d, reason: %s\n", w.attempts, err.Error())
	}

	return err
}

func (w *WebTranslateIt) update() error {
	w.attempts = w.attempts + 1

	for i := range w.dictionaries {
		w.dictionaries[i].Update = false
	}

	project, err := w.wti.GetProject()
	if err != nil {
		return err
	}

	if project.UpdatedAt == w.updateAt {
		return nil
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

	w.updateAt = project.UpdatedAt
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

	return nil
}

func (w *WebTranslateIt) parseFile(file wti.File, content []byte) {
	dictionary, err := w.GetDictionary(file.LocaleCode)
	if err != nil {
		dictionary = &Dictionary{
			Locale: GetLocale(file.LocaleCode),
		}

		w.dictionaries = append(w.dictionaries, dictionary)
	}

	dictionary.File = file
	if dictionary.Hash != file.Hash {
		dictionary.Phrases = map[string]string{}
		dictionary.Update = true
		dictionary.Hash = file.Hash

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
