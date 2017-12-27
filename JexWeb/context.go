package jexweb

import (
	"reflect"

	"github.com/JexLib/golang/JexWeb/session"
	"github.com/labstack/echo"
)

type Context interface {
	echo.Context
	_init(echo.Context)
	Init()
	SetLayout(layout string)
	Session() session.Session
	// SetLayout(string)
	// Set()
	// Get()
	// Render()
}

type JexContext struct {
	echo.Context
	Data map[string]interface{}
}

func NewContext(ctx echo.Context) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := reflect.TypeOf(ctx)
			e := t.Elem()
			var v = reflect.New(e)
			jc := v.Interface().(Context)
			jc._init(c)
			jc.Init()
			return h(jc)
		}
	}
}

func (c *JexContext) _init(ctx echo.Context) {
	c.Context = ctx
	c.Data = map[string]interface{}{}
	// c.Flash = Flash{c}

}

func (c *JexContext) Init() {
}

func (c *JexContext) SetLayout(layout string) {
	c.Context.Set("Layout", layout)
}

func (c *JexContext) Set(key string, val interface{}) {
	c.Data[key] = val
}

func (c *JexContext) Get(key string) interface{} {
	return c.Data[key]
}

func (c *JexContext) Session() session.Session {
	return session.Default(c.Context)
}
