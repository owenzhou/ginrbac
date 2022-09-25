package app

import (
	"strconv"

	"github.com/owenzhou/ginrbac/support/facades"

	"github.com/go-playground/validator/v10"
)

//将validator错误信息翻译成中文
func Translate(err error) []string {
	var result = make([]string, 0)
	switch err.(type) {
	case *strconv.NumError:
		errors := err.Error()
		result = append(result, errors)
		return result
	default:
		errors := err.(validator.ValidationErrors)
		for _, err := range errors {
			result = append(result, err.Translate(facades.Translator))
		}
		return result
	}
}
