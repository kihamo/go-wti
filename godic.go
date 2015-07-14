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

	handler := turnpike.NewBasicWebsocketServer(Realm)
	client, err := handler.GetLocalClient(Realm)
	if err != nil {
		return nil, err
	}

	if err = client.Register(GetDictionaryMethod, server.getDictionary); err != nil {
		return nil, err
	}

	if err = client.Register(DictionaryUpdateMethod, server.dictionaryUpdate); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/dictionaries", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("dictionaries health check"))
	})
	mux.Handle("/", handler)

	return &GodicServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		client: client,
	}, nil
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
	s.dictionaryUpdate(nil, nil)
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
