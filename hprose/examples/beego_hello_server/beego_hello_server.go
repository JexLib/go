package main

import (
	"github.com/JexLib/golang/hprose/rpc"
	"github.com/astaxie/beego"
)

func hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	service := rpc.NewHTTPService()
	service.AddFunction("hello", hello)
	beego.Handler("/hello", service)
	beego.Run()
}
