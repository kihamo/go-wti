package gowti

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	wti "github.com/fromYukki/webtranslateit_go_client"
	"github.com/kihamo/go-wti/gen-go/translator"
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

const (
	updatedLangRegexp = "(?is)\\$lang\\[\"(?P<key>.*?)\"\\] = \"(?P<value>.*?)\";"
)

func init() {
	var err error

	exp, err = regexp.Compile(updatedLangRegexp)
	if err != nil {
		panic(err)
	}
}

type TranslatorHandler struct {
	wti                 *wti.WebTranslateIt
	updateRetryDelay    time.Duration
	updateRetryAttempts int64
	dictionaries        map[string]map[string]string
}

func NewTranslatorHandler(wtiToken string, updateRetryDelay time.Duration, updateRetryAttempts int64) *TranslatorHandler {
	h := &TranslatorHandler{
		wti:                 wti.NewWebTranslateIt(wtiToken),
		updateRetryDelay:    updateRetryDelay,
		updateRetryAttempts: updateRetryAttempts,
		dictionaries:        map[string]map[string]string{},
	}
	h.run()

	return h
}

func (h *TranslatorHandler) Ping() (bool, error) {
	return true, nil
}

func (h *TranslatorHandler) GetDictionary(locale string) (map[string]string, error) {
	locale = strings.ToLower(locale)

	dictionary, ok := h.dictionaries[locale]
	if !ok {
		return nil, &translator.TranslatorError{
			ErrorCode:    translator.TranslatorErrorCode_LOCALE_NOT_FOUND,
			ErrorMessage: fmt.Sprintf("Locale %s nof found", locale),
		}
	}

	return dictionary, nil
}

func (h *TranslatorHandler) run() {
	project, err := h.wti.GetProject()
	if err != nil {
		// TODO: log
		return
	}

	zipFile, err := project.ZipFile()
	if err != nil {
		// TODO: log
		return
	}

	data, err := zipFile.Extract()
	if err != nil {
		// TODO: log
		return
	}

	for i := range project.ProjectFiles {
		h.parseFile(project.ProjectFiles[i], data[project.ProjectFiles[i].Name])
	}

	if h.updateRetryDelay > 0 {
		time.AfterFunc(h.updateRetryDelay, func() { h.run() })
	}
}

func (h *TranslatorHandler) parseFile(file wti.File, content []byte) {
	locale := file.LocaleCode
	if _, ok := localAliases[locale]; ok {
		locale = localAliases[locale]
	}

	locale = strings.ToLower(locale)
	h.dictionaries[locale] = map[string]string{}

	for _, match := range exp.FindAllStringSubmatch(string(content), -1) {
		value := match[2]
		if value == "" {
			value = match[1]
		}

		h.dictionaries[locale][match[1]] = value
	}

	log.Printf("Update %s locale", locale)
}
