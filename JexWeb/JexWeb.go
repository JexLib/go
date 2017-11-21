package jexweb

import (
	"fmt"
	"html/template"
	"reflect"

	"golang.org/x/crypto/acme/autocert"

	"github.com/JexLib/golang/JexWeb/session"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/render"
)

const (
	banner_jex = `
       _________  __
      / / ____/ |/ /
 __  / / __/  |   /
/ /_/ / /___ /   |
\____/_____//_/|_|  
`
	banner_jexweb = `
       _________  __              __
      / / ____/ |/ /      _____  / /_
 __  / / __/  |   / | /| / / _ \/ __ \
/ /_/ / /___ /   || |/ |/ /  __/ /_/ /
\____/_____//_/|_||__/|__/\___/_.___/
`
)

type (
	JexWeb struct {
		Config Config
		Echo   *echo.Echo
		// Perm          *permissionbolt.Permissions
		// _denyFunction echo.HandlerFunc
		//controllers   map[string]iController
	}

	// HandlerFunc func() error
)

func defConfig() Config {
	return Config{
		HttpPort:      8080,
		AssetsDir:     "public/assets",
		PublicDir:     "public",
		TemplateDir:   "templates",
		AppLayout:     "layout",
		IsDevelopment: true,
	}
}

func New(store ...session.Store) *JexWeb {
	return NewWithConfig(defConfig(), store...)
}

func NewWithConfig(config Config, store ...session.Store) *JexWeb {
	jwb := &JexWeb{
		Config: config,
		Echo:   echo.New(),
		// _denyFunction: permissionDenied,
		//	controllers:        make(map[string]iController),
	}
	jwb.Echo.Use(middleware.Recover())

	// jwb.Echo.Use(middleware.Logger())
	jwb.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	if len(store) > 0 {
		jwb.Echo.Use(session.Sessions("SESSID", store[0]))
	}
	// jwb.Perm, _ = permissionbolt.NewWithConf("permdb")
	// jwb.Perm.UserState().SetCookieTimeout(60 * 60)
	// store := session.NewFileSystemStoreStore("store")
	// jwb.Echo.Use(session.Sessions("SESSID", store))
	return jwb
}

// func (jwb *JexWeb) UsePermissionMW(beforeMiddleware ...echo.MiddlewareFunc) {
// 	jwb.Echo.Use(beforeMiddleware...)
// 	jwb.Echo.Use(jex_middleware.PermissionMiddleware(jwb.Perm, jwb.denyFunction))
// }

// func permissionDenied(c echo.Context) error {
// 	c.Error(echo.ErrForbidden)
// 	return nil
// 	//return c.String(http.StatusForbidden, "Permission denied!")
// }

// func (jwb *JexWeb) denyFunction(c echo.Context) error {
// 	// if web._denyFunction == nil {
// 	// 	return c.String(http.StatusForbidden, "Permission denied!")
// 	// } else {
// 	return jwb._denyFunction(c)
// 	// }
// }

// func (jwb *JexWeb) SetDenyFunction(denyFunction echo.HandlerFunc) {
// 	jwb._denyFunction = denyFunction
// }

// func (jwb *JexWeb) Group(prefix string, controller iController, m ...echo.MiddlewareFunc) *Group {

// 	return &Group{
// 		controller: controller,
// 		Group:      jwb.Echo.Group(prefix, m...),
// 	}
// }

func defHandleFunc(handlerFuncName string, controller iController) echo.HandlerFunc {
	return func(c echo.Context) error {
		t := reflect.TypeOf(controller)
		e := t.Elem()
		var v = reflect.New(e)
		jc := v.Interface().(iController)
		jc.init(c)
		jc.Init()
		rets := v.MethodByName(handlerFuncName).Call([]reflect.Value{})
		if len(rets) > 0 {
			if err, ok := rets[0].Interface().(error); ok {
				return err
			}
		}

		return nil
	}
}

type HandlerFunc func(*JexContext) error

func (jwb *JexWeb) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	fmt.Println("path:", path)
	jwb.Echo.GET(path, func(c echo.Context) error {
		t := reflect.TypeOf(h)

		fmt.Printf("参数数量:", t.NumIn())
		fmt.Println("参数类型", t.In(0))
		i := 0
		for ; i < t.NumIn()-1; i++ {
			fmt.Println("    ┣", t.In(i)) // 获取参数类型
		}

		// fmt.Printf("\n%-8v %v 个方法:\n", v, v.NumMethod())
		// fmt.Println("ttt:",)
		for i := 0; i < t.NumMethod(); i++ {
			fmt.Println("Method:", t.Method(i).Name)
		}

		// t := reflect.TypeOf(h)
		// e := t.Elem()
		// var v = reflect.New(e)
		// jc := v.Interface().(iController)
		// jc.init(c)
		// jc.Init()

		return nil
	}, m...)
}

// func (jwb *JexWeb) GET(path string, controller iController, handlerFuncName string, m ...echo.MiddlewareFunc) {
// 	jwb.Echo.GET(path, defHandleFunc(handlerFuncName, controller), m...)
// }

// func (g *Group) GET(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.GET(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) POST(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.POST(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) Any(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.Any(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) DELETE(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.DELETE(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) OPTIONS(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.OPTIONS(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) PATCH(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.PATCH(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) PUT(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.PUT(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) HEAD(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.HEAD(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) CONNECT(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.CONNECT(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) TRACE(path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.TRACE(path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (g *Group) Match(methods []string, path string, handlerFuncName string, m ...echo.MiddlewareFunc) *Group {
// 	g.Group.Match(methods, path, defHandleFunc(handlerFuncName, g.controller), m...)
// 	return g
// }

// func (jwb *JexWeb) Any(path string, controller iController, handler string, m ...echo.MiddlewareFunc) {
// 	jwb.Echo.Any(path, jwb.routeHandleName(controller, handler), m...)
// }

// func (jwb *JexWeb) GET(path string, controller iController, handler string, m ...echo.MiddlewareFunc) {
// 	jwb.Echo.GET(path, jwb.routeHandleName(controller, handler), m...)
// }

// func (jwb *JexWeb) Group(prefix string, controller iController, routes []Route, m ...echo.MiddlewareFunc) {
// 	g := jwb.Echo.Group(prefix, m...)
// 	//	jwb.controllers[prefix] = controller

// 	for _, v := range routes {
// 		switch v.Method {
// 		case "Any":
// 			g.Any(v.Path, jwb.routeHandleName(controller, v.Handler))
// 		case "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE", "SOCKET":
// 			jwb.Echo.Match([]string{v.Method}, v.Path, jwb.routeHandleName(controller, v.Handler))
// 		}
// 	}
// 	//增加handle
// 	//`roud:""`

// 	// jwb.Echo.Match
// 	//st := reflect.TypeOf(controller)
// 	//field := st.Field(0)
// 	//    fmt.Println(field.Tag.Get("color"), field.Tag.Get("species")
// 	//jwb._denyFunction = denyFunction
// }

// func (jwb *JexWeb) routeHandleName(controller iController, handel string) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		t := reflect.TypeOf(controller)
// 		e := t.Elem()
// 		var v = reflect.New(e)
// 		v.Interface().(iController).Init(c)
// 		rets := v.MethodByName(handel).Call([]reflect.Value{})
// 		if len(rets) > 0 {
// 			if err, ok := rets[0].Interface().(error); ok {
// 				return err
// 			}
// 		}

// 		return nil
// 	}
// }

func (jwb *JexWeb) Start() {
	_Render := render.New(render.Options{
		Directory:                 jwb.Config.TemplateDir,
		Asset:                     nil,
		AssetNames:                nil,
		Layout:                    jwb.Config.AppLayout,
		Extensions:                []string{".html", ".tmpl"},
		Funcs:                     []template.FuncMap{_helperFuncs, _ExtendFuncs},
		Delims:                    render.Delims{"{{", "}}"},
		Charset:                   "UTF-8",
		IndentJSON:                false,
		IndentXML:                 false,
		PrefixJSON:                []byte(""),
		PrefixXML:                 []byte(""),
		HTMLContentType:           "text/html",
		IsDevelopment:             true,
		UnEscapeHTML:              false,
		StreamingJSON:             false,
		RequirePartials:           false,
		DisableHTTPErrorRendering: false,
	})
	r := &RenderWrapper{_Render}
	jwb.Echo.Renderer = r

	jwb.Echo.Binder = &binder{}
	jwb.Echo.Static("assets", jwb.Config.AssetsDir)
	jwb.Echo.Static("public", jwb.Config.PublicDir)
	jwb.Echo.HideBanner = true
	fmt.Println(banner_jexweb)
	if jwb.Config.Https.Enabled {
		jwb.Echo.Pre(middleware.HTTPSRedirect())
		jwb.Echo.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		address := fmt.Sprintf("%s:%d", jwb.Config.Addr, jwb.Config.Https.Port)
		fmt.Println("Starting JexWeb listening at", address)
		jwb.Echo.Logger.Fatal(jwb.Echo.StartAutoTLS(address))
	} else {
		address := fmt.Sprintf("%s:%d", jwb.Config.Addr, jwb.Config.HttpPort)
		fmt.Println("Starting JexWeb listening at", address)
		jwb.Echo.Logger.Fatal(jwb.Echo.Start(address))
	}
}
