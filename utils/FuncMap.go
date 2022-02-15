package utils

import (
	"github.com/owenzhou/ginrbac/utils/php"
	"html/template"
)

//传给模板的函数
var FuncMap = template.FuncMap{
	//时间处理
	"date": php.Date,
	//判断字符串是否存在
	"strstr": php.Strstr,
	//字符串截取
	"substr": php.Substr,
	//字符串中文截取
	"mb_substr": php.Mb_substr,
	//字符串长度
	"strlen": php.Strlen,
	//中文字符串长度
	"mb_strlen": php.Mb_strlen,
	//字符串反转
	"strrev": php.Strrev,
	//字符替换
	"str_replace": php.Str_replace,
	//生成url
	"url":       URL,
	"range_int": RangeInt,
	"trim":      php.Trim,
	//将html实体转为html
	"parseHtml": ParseHtml,
}
