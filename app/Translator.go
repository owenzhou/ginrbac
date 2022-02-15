package app

import (
	"ginrbac/bootstrap/support/facades"

	"github.com/go-playground/validator/v10"
)

//将validator错误信息翻译成中文
func Translate(err error) []string {
	var result = make([]string, 0)
	errors := err.(validator.ValidationErrors)
	for _, err := range errors {
		result = append(result, err.Translate(facades.Translator))
	}
	return result
}
