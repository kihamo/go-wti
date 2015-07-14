package godic // import "github.com/kihamo/godic"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/jcelliott/turnpike.v2"
)

const (
	Realm                  = "godic"
	UpdateTopic            = "dictionary.updater.%s"
	DictionaryUpdateMethod = "dictionary.update"
	GetDictionaryMethod    = "dictionary.get"
	GetDictionaryErrorURI  = "dictionary.get.error"
)

type GodicServer struct {
	server *http.Server
	client *turnpike.Client
	wti    *WebTranslateIt
}

func NewServer(addr string, wti *WebTranslateIt, debug bool) (server *GodicServer, err error) {
	if debug {
		turnpike.Debug()
	}

	handler := turnpike.NewBasicWebsocketServer(Realm)
	client, err := handler.GetLocalClient(Realm)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	server = &GodicServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		client: client,
		wti:    wti,
	}

	if err = client.Register(GetDictionaryMethod, server.GetDictionary); err != nil {
		return nil, err
	}

	if err = client.Register(DictionaryUpdateMethod, server.DictionaryUpdate); err != nil {
		return nil, err
	}

	mux.HandleFunc("/dictionaries", func(w http.ResponseWriter, r *http.Request) {
		dictionaries := server.wti.GetDictionaries()
		reply := make(map[string]map[string]interface{}, len(dictionaries))

		for i := range dictionaries {
			reply[dictionaries[i].Locale] = map[string]interface{}{
				"count":     len(dictionaries[i].Phrases),
				"hash":      dictionaries[i].File.Hash,
				"update_at": dictionaries[i].File.UpdatedAt,
			}
		}

		response, _ := json.Marshal(reply)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(response)
	})
	mux.Handle("/", handler)

	wti.SetCallback(func(locales []string) {
		for i := range locales {
			client.Publish(fmt.Sprintf(UpdateTopic, locales[i]), nil, nil)
		}
	})

	return server, nil
}

func (s *GodicServer) ListenAndServe() error {
	s.DictionaryUpdate(nil, nil)
	return s.server.ListenAndServe()
}

func (s *GodicServer) GetDictionary(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	if len(args) < 1 {
		return &turnpike.CallResult{Err: turnpike.URI(GetDictionaryErrorURI)}
	}

	locale, ok := args[0].(string)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI(GetDictionaryErrorURI)}
	}

	dictionary, err := s.wti.GetDictionary(locale)
	if err != nil {
		return &turnpike.CallResult{Err: turnpike.URI(GetDictionaryErrorURI)}
	}

	return &turnpike.CallResult{Args: []interface{}{dictionary.Phrases}}
}

func (s *GodicServer) DictionaryUpdate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	go func() {
		err := s.wti.Update()
		if err != nil {
			log.Printf("Update dictionaries error %s\n", err)
		}
	}()

	return &turnpike.CallResult{Args: []interface{}{true}}
}
