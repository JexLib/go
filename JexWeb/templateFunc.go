package jexweb

import (
	"html/template"
	"strings"
)

var (
	_helperFuncs template.FuncMap
	_ExtendFuncs template.FuncMap
)

func AddTemplateFunc(key string, Func interface{}) {
	_ExtendFuncs[key] = Func
}

func In(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

//字符串转大写
func ToUpper(val string) string {
	return strings.ToUpper(val)
}

//字符串转小写
func ToLower(val string) string {
	return strings.ToLower(val)
}

func init() {
	_helperFuncs = make(template.FuncMap)
	_ExtendFuncs = make(template.FuncMap)

	AddTemplateFunc("jex_In", In)
	AddTemplateFunc("jex_ToUpper", ToUpper)
	AddTemplateFunc("jex_ToLower", ToLower)
}
