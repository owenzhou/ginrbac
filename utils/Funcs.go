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
	"regexp"
	"sort"
	"strconv"
	"strings"

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

//当前url添加参数
func UrlAppend(k string, v interface{}, currentUrl string) string {
	var reg *regexp.Regexp
	if strings.Contains(currentUrl, "?") {
		reg = regexp.MustCompile(`([\?&]{1}` + k + `=)([^&]*)`)
	} else {
		reg = regexp.MustCompile(`([/]{1}` + k + `/)([^/.]*)`)
	}
	urlArr := strings.Split(currentUrl, ".")
	var suffix = ""
	if len(urlArr) > 1{
		suffix = "." + urlArr[1]
	}
	if !reg.MatchString(currentUrl) {
		if strings.Contains(currentUrl, "?"){
			query, _ := url.ParseQuery(currentUrl)
			query.Add(k, Int2Str(v))
			u, _ := url.PathUnescape(query.Encode())
			return u
		}
		return urlArr[0] + "/" + k + "/" + Int2Str(v) + suffix
	}
	return reg.ReplaceAllString(currentUrl, "${1}" + Int2Str(v))
}

//当前url删除参数
func UrlDelete(k, currentUrl string) string{
	var reg *regexp.Regexp
	if strings.Contains(currentUrl, "?") {
		reg = regexp.MustCompile(`([\?&]{1}` + k + `=)([^&]*)`)
	} else {
		reg = regexp.MustCompile(`([/]{1}` + k + `/)([^/.]*)`)
	}

	if reg.MatchString(currentUrl) {
		return reg.ReplaceAllString(currentUrl, "")
	}
	return currentUrl
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

func ParseJS(str string) template.JS{
	return template.JS(str)
}

//可填参数data,只取第一个
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

	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}

	b := bytes.NewBuffer([]byte{})
	err1 := tmpl.ExecuteTemplate(b, tname, d)
	if err1 != nil {
		panic(err1)
	}
	
	return template.HTML(b.String())
}

//判断指针是否相等
func DeepEqual(val, val2 interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(val))
	v2 := reflect.Indirect(reflect.ValueOf(val2))
	var newV, newV2 interface{}
	
	switch(v.Kind()){
		//Int、Int8、Int16、Int32、Int64
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newV = v.Int()
		//Uint、Uintptr、Uint8、Uint16、Uint32、Uint64
		case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			newV = int64(v.Uint())
		case reflect.Float32, reflect.Float64:
			newV = v.Float()
		case reflect.String:
			newV = v.String()
	}

	switch(v2.Kind()){
		//Int、Int8、Int16、Int32、Int64
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newV2 = v2.Int()
		//Uint、Uintptr、Uint8、Uint16、Uint32、Uint64
		case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			newV2 = int64(v2.Uint())
		case reflect.Float32, reflect.Float64:
			newV2 = v2.Float()
		case reflect.String:
			newV2 = v2.String()
	}
	
	return newV == newV2
}

func DeepNotEqual(val, val2 interface{}) bool {
	return !DeepEqual(val, val2)
}

func PlusInt(v, i int) int {
	return v + i
}

//获取指针地址保存的值
func PtrValue(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return v
	}
	return reflect.Indirect(val).Interface()
}
