package main

import (
	"net/http"

	"github.com/JexLib/golang/hprose/rpc"
)

func hello(name string, context *rpc.HTTPContext) string {
	return "Hello !" + context.Request.RemoteAddr
}

func main() {
	service := rpc.NewHTTPService()
	service.Debug = true
	service.AddFunction("hello", hello)
	http.ListenAndServe(":8080", service)
}
