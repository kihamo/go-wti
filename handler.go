package gowti

type TranslatorHandler struct {
}

func NewTranslatorHandler() *TranslatorHandler {
	return &TranslatorHandler{}
}

func (h *TranslatorHandler) Ping() (bool, error) {
	return true, nil
}
