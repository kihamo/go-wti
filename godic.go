package godic // import "github.com/kihamo/godic"

import (
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

func NewServer(addr string, debug bool) (server *GodicServer, err error) {
	if debug {
		turnpike.Debug()
	}

	s := turnpike.NewBasicWebsocketServer(Realm)
	server = &GodicServer{
		server: &http.Server{
			Handler: s,
			Addr:    addr,
		},
	}

	// TODO: handle /dictionaries

	server.client, err = s.GetLocalClient(Realm)
	if err != nil {
		return nil, err
	}

	if err = server.client.Register(GetDictionaryMethod, server.getDictionary); err != nil {
		return nil, err
	}

	if err = server.client.Register(DictionaryUpdateMethod, server.dictionaryUpdate); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *GodicServer) SetWebTranslateIt(wti *WebTranslateIt) {
	s.wti = wti
	s.wti.SetCallback(func(locales []string) {
		for i := range locales {
			s.client.Publish(fmt.Sprintf(UpdateTopic, locales[i]), nil, nil)
		}
	})
}

func (s *GodicServer) ListenAndServe() error {
	s.client.Call(DictionaryUpdateMethod, nil, nil)
	return s.server.ListenAndServe()
}

func (s *GodicServer) getDictionary(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
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

func (s *GodicServer) dictionaryUpdate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	go func() {
		err := s.wti.Update()
		if err != nil {
			log.Printf("Update dictionaries error %s\n", err)
		}
	}()

	return &turnpike.CallResult{Args: []interface{}{true}}
}
