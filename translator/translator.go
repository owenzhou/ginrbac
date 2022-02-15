package translator

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
)

func newTranslator() ut.Translator {
	zh := zh.New()
	utr := ut.New(zh)
	trans, _ := utr.GetTranslator("zh")
	validate := binding.Validator.Engine().(*validator.Validate)
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("label")
		return name
	})
	zh_translation.RegisterDefaultTranslations(validate, trans)
	return trans
}
