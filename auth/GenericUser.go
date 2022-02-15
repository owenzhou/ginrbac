package auth

import (
	"reflect"
	"strconv"
)

type GenericUser struct {
	Attributes map[string]interface{}
}

func (g *GenericUser) Get(name string) interface{} {
	if g.Attributes == nil {
		return ""
	}

	if v, ok := g.Attributes[name]; ok {
		return v
	}
	return ""
}

func (g *GenericUser) GetAuthIdentifierName() string {
	return "id"
}

func (g *GenericUser) GetAuthIdentifier() string {
	id := g.Attributes[g.GetAuthIdentifierName()]
	t := reflect.TypeOf(id)
	if t.Kind() == reflect.Uint32 {
		return strconv.FormatUint(uint64(id.(uint32)), 10)
	}
	if t.Kind() == reflect.Uint64 {
		return strconv.FormatUint(id.(uint64), 10)
	}
	if t.Kind() == reflect.Int32 {
		return strconv.FormatInt(int64(id.(int32)), 10)
	}
	if t.Kind() == reflect.Int64 {
		return strconv.FormatInt(id.(int64), 10)
	}
	return id.(string)
}

func (g *GenericUser) GetAuthPassword() string {
	return g.Attributes["password"].(string)
}

func (g *GenericUser) SetRememberToken(v interface{}) {
	g.Attributes[g.GetRememberTokenName()] = v
}

func (g *GenericUser) GetRememberToken() string {
	return g.Attributes[g.GetRememberTokenName()].(string)
}

func (g *GenericUser) GetRememberTokenName() string {
	return "remember_token"
}

func (g *GenericUser) SetApiToken(v interface{}) {
	g.Attributes[g.GetApiTokenName()] = v
}

func (g *GenericUser) GetApiToken() string {
	return g.Attributes[g.GetApiTokenName()].(string)
}

func (g *GenericUser) GetApiTokenName() string {
	return "api_token"
}

func (g *GenericUser) IsEmpty() bool {
	return len(g.Attributes) <= 0
}
