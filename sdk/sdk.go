package sdk

import (
	"errors"
	"fmt"

	"github.com/kihamo/godic"
	"gopkg.in/jcelliott/turnpike.v2"
)

// TODO: reconnect after failed request

type GodicClient struct {
	client *turnpike.Client
}

func NewClient(addr string, debug bool) (*GodicClient, error) {
	c, err := turnpike.NewWebsocketClient(turnpike.MSGPACK, fmt.Sprintf("ws://%s/", addr))
	if err != nil {
		return nil, err
	}

	_, err = c.JoinRealm(godic.Realm, turnpike.SUBSCRIBER|turnpike.CALLER, nil)
	if err != nil {
		return nil, err
	}

	return &GodicClient{
		client: c,
	}, nil
}

func (s *GodicClient) GetDictionary(locale string) (map[string]string, error) {
	response, err := s.client.Call(godic.GetDictionaryMethod, []interface{}{locale}, nil)
	if err != nil {
		return nil, err
	}

	dic, ok := response.Arguments[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("Error type")
	}

	result := make(map[string]string, len(dic))
	for in, out := range dic {
		result[in] = out.(string)
	}

	return result, nil
}

func (s *GodicClient) DictionaryUpdate() (bool, error) {
	response, err := s.client.Call(godic.DictionaryUpdateMethod, nil, nil)
	if err != nil {
		return false, err
	}

	result, ok := response.Arguments[0].(bool)
	if !ok {
		return false, errors.New("Error type")
	}

	return result, nil
}

func (s *GodicClient) UpdateSubscribe(locales []string, f func(locale string)) {
	for i := range locales {
		locale := godic.GetLocale(locales[i])

		topic := fmt.Sprintf(godic.UpdateTopic, locale)
		s.client.Subscribe(topic, func([]interface{}, map[string]interface{}) {
			f(locale)
		})
	}
}
