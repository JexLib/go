package jexweb

import "html/template"

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

func init() {
	_helperFuncs = make(template.FuncMap)
	_ExtendFuncs = make(template.FuncMap)

	AddTemplateFunc("jex_in", In)
}
