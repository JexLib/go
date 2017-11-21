package jexweb

import (
	"github.com/labstack/echo"
)

type iController interface {
	init(echo.Context)
	Init()
	SetLayout(layout string)

	// Set(key string, val interface{})
	// Get(key string) interface{}
	// Data() map[string]interface{}

	// SetLayout(string)
	// Set()
	// Get()
	// Render()
}

// var (
// 	permissions *permissionbolt.Permissions
// )

type Controller struct {
	echo.Context
	Data  map[string]interface{}
	Flash Flash
	user  *User
}

type Flash struct {
	controller *Controller
}

type User struct {
	Name string
}

func (c *Controller) init(ctx echo.Context) {
	c.Context = ctx
	c.Data = map[string]interface{}{}
	c.Flash = Flash{c}
	//	_helperFuncs["Flashs"] = c.Flash.Data()
}

func (c *Controller) Init() {

}

func (c *Controller) SetLayout(layout string) {
	c.Data["@LAYOUT@"] = layout

}

func (c *Controller) SetData(key string, val interface{}) {
	c.Data[key] = val
}

func (c *Controller) GetData(key string) interface{} {
	return c.Data[key]
}

// func (c *Controller) Session() session.Session {
// 	return session.Default(c.Context)
// }

// func (f *Flash) Set(key string, value interface{}) {
// 	f.controller.Session().AddFlash(map[string]interface{}{
// 		key: value,
// 	})
// }

// func (f *Flash) Info(value interface{}) {
// 	f.Set("info", value)
// }

// func (f *Flash) Error(value interface{}) {
// 	f.Set("error", value)
// }

// func (f *Flash) Success(value interface{}) {
// 	f.Set("success", value)
// }

// func (f *Flash) Save() {
// 	f.controller.Session().Save()
// }

// func (f *Flash) Data() []interface{} {
// 	return f.controller.Session().Flashes()
// }

// func (c *Controller) Render(code int, name string, data interface{}) error {

// }

// func (r *Render) HTML(w io.Writer, status int, name string, binding interface{}, htmlOpt ...HTMLOptions) error {
