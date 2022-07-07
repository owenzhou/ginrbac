package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"html/template"
	"net/url"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"

	"github.com/owenzhou/ginrbac/support/facades"
)

//int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 float32 float64转字符串
func Int2Str(i interface{}) string {
	var v string
	switch n := i.(type) {
	case int:
		v = strconv.Itoa(n)
	case int8:
		v = strconv.FormatInt(int64(n), 10)
	case int16:
		v = strconv.FormatInt(int64(n), 10)
	case int32:
		v = strconv.FormatInt(int64(n), 10)
	case int64:
		v = strconv.FormatInt(n, 10)
	case uint:
		v = strconv.FormatUint(uint64(n), 10)
	case uint8:
		v = strconv.FormatUint(uint64(n), 10)
	case uint16:
		v = strconv.FormatUint(uint64(n), 10)
	case uint32:
		v = strconv.FormatUint(uint64(n), 10)
	case uint64:
		v = strconv.FormatUint(n, 10)
	case float32:
		v = fmt.Sprintf("%f", n)
	case float64:
		v = fmt.Sprintf("%f", n)
	case string:
		v = n
	default:
		v = ""
	}
	return v
}

//字符串转int
func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

//生成带query的url
func URL(route string, params ...interface{}) string {
	if len(params)%2 == 1 {
		return route
	}
	values := url.Values{}
	if len(params) == 2 {
		values.Add(params[0].(string), Int2Str(params[1]))
	} else {
		for i := 0; i < len(params); i += 2 {
			values.Add(params[i].(string), Int2Str(params[i+1]))
		}
	}
	p, _ := url.PathUnescape(values.Encode())
	return route + "?" + p
}

//生成int数组
func RangeInt(start, end int) []int {
	var startTmp = start
	var isSort = false
	if start > end {
		startTmp = end
		end = start
		start = startTmp
		isSort = true
	}
	m := make([]int, 0)
	for i := start; i <= end; i++ {
		m = append(m, i)
	}
	if isSort {
		sort.Sort(sort.Reverse(sort.IntSlice(m)))
	}
	return m
}

func Struct2Map(v interface{}) map[string]interface{} {
	mapData := make(map[string]interface{})
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		panic("Struct2Map error: param is not a struct")
	}
	b, _ := json.Marshal(v)

	json.Unmarshal(b, &mapData)
	return mapData
}

func ParseHtml(str string) template.HTML {
	return template.HTML(str)
}

func Widget(name string, data ...interface{}) template.HTML {
	widget, err := fs.Glob(facades.Views, "views/" + name + "*")
	if err != nil {
		panic(err)
	}
	
	tname := filepath.Base(widget[0])
	funcMap := template.FuncMap{}
	for n, f := range FuncMap{
		funcMap[n] = f
	}
	funcMap["widget"] = Widget

	var tmpl *template.Template
	if facades.Config.Debug{
		tmpl = template.Must(template.New(tname).Funcs(funcMap).ParseFiles(widget...))
	}else{
		tmpl = template.Must(template.New(tname).Funcs(funcMap).ParseFS(facades.Views, widget...))
	}

	b := bytes.NewBuffer([]byte{})
	err1 := tmpl.ExecuteTemplate(b, tname, data)
	if err1 != nil {
		panic(err1)
	}
	
	return template.HTML(b.String())
}
