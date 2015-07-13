package godic

import (
	"errors"
	"regexp"
	"strings"
	"time"

	wti "github.com/fromYukki/webtranslateit_go_client"
)

const (
	updatedLangRegexp = "(?is)\\$lang\\[\"(?P<key>.*?)\"\\] = \"(?P<value>.*?)\";"
)

var (
	exp          *regexp.Regexp
	localAliases = map[string]string{
		"en_en_VN": "en_VN",
		"en_en_TH": "en_TH",
		"en_en_ID": "en_ID",
		"ms":       "ms_MY",
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
	dictionaries        map[string]map[string]string
	callback            func()
}

func NewWebTranslateIt(wtiToken string, updateRetryDelay time.Duration, updateRetryAttempts int64) *WebTranslateIt {
	return &WebTranslateIt{
		wti:                 wti.NewWebTranslateIt(wtiToken),
		updateRetryDelay:    updateRetryDelay,
		updateRetryAttempts: updateRetryAttempts,
		dictionaries:        map[string]map[string]string{},
	}
}

func (w *WebTranslateIt) GetDictionary(locale string) (map[string]string, error) {
	locale = strings.ToLower(locale)

	dictionary, ok := w.dictionaries[locale]
	if !ok {
		return nil, errors.New("Locale not exists")
	}

	return dictionary, nil
}

func (w *WebTranslateIt) SetCallback(f func()) {
	w.callback = f
}

func (w *WebTranslateIt) Update() error {
	project, err := w.wti.GetProject()
	if err != nil {
		// TODO: log
		return err
	}

	zipFile, err := project.ZipFile()
	if err != nil {
		// TODO: log
		return err
	}

	data, err := zipFile.Extract()
	if err != nil {
		// TODO: log
		return err
	}

	for i := range project.ProjectFiles {
		w.parseFile(project.ProjectFiles[i], data[project.ProjectFiles[i].Name])
	}

	if w.callback != nil {
		w.callback()
	}

	if w.updateRetryDelay > 0 {
		time.AfterFunc(w.updateRetryDelay, func() { w.Update() })
	}

	return nil
}

func (w *WebTranslateIt) parseFile(file wti.File, content []byte) {
	locale := file.LocaleCode
	if _, ok := localAliases[locale]; ok {
		locale = localAliases[locale]
	}

	locale = strings.ToLower(locale)
	w.dictionaries[locale] = map[string]string{}

	for _, match := range exp.FindAllStringSubmatch(string(content), -1) {
		value := match[2]
		if value == "" {
			value = match[1]
		}

		w.dictionaries[locale][match[1]] = value
	}
}
