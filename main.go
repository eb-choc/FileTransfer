package main

import (
	"github.com/kataras/iris"
	"iris.projects/test1/router"
)

func main() {
	app := iris.Default()
	tpl := iris.HTML("./views", ".html").Reload(true)
	router.RegisterTplFuns(tpl)
	app.RegisterView(tpl)
	app.HandleDir("/static/", "./static/")
	router.InitRouter(app)
	app.Run(iris.Addr(":90"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:         false,
		DisableAutoFireStatusCode: false,
		DisableInterruptHandler:   false,
	}))
}
