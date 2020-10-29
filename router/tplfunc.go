package router

import (
	"github.com/kataras/iris/view"
	"strings"
)

func RegisterTplFuns(v *view.HTMLEngine) {
	v.AddFunc("replaceStrPlus", ReplacePlus)
}


func ReplacePlus(str string) string {
	return strings.Replace(str, "&#43;", "+", -1);
}