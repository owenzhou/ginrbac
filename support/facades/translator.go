package facades

import (
	ut "github.com/go-playground/universal-translator"
)

var Translator ut.Translator

type TranslatorFacade struct {
	*Facade
}

func (t *TranslatorFacade) GetFacadeAccessor() {
	Translator = t.App.Make("translator").(ut.Translator)
}
