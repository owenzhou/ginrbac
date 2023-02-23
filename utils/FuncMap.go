package utils

import (
	"html/template"

	"github.com/owenzhou/ginrbac/utils/php"
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
	//当前url添加参数
	"url_append": UrlAppend,
	//当前url删除参数
	"url_delete": UrlDelete,
	"range_int": RangeInt,
	"trim":      php.Trim,
	//将html实体转为html
	"parseHtml": ParseHtml,
	"parseJS": ParseJS,
	//字符串转数组
	"split": php.Split,
	//判断指针值是否相等
	"deep_equal": DeepEqual,
	"deep_notequal": DeepNotEqual,
	//将各种int转为str进行比较
	"int2str": Int2Str,
	"md5": php.Md5,
	"strip_tags": php.Strip_tags,
	"ptr_value": PtrValue,
	"plus_int": Plus[int],
	"minus_int": Minus[int],
	"multiply_int": Multiply[int],
	"divide_int": Divide[int],
	"remainder_int": Remainder[int],
}
