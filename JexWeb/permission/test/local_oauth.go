package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JexLib/golang/JexWeb/session"

	"github.com/JexLib/golang/JexWeb/permission"
	"github.com/JexLib/golang/cache/memory"
	"github.com/labstack/echo"
)

func OnLoginBefore(account, password string, us *permission.User) error {
	if account == "admin" && password == "112233" {
		us.IsAdmin = true
		return nil
	}
	if account == "user" && password == "112233" {
		return nil
	}

	return errors.New("验证用户失败")
}

func main() {
	e := echo.New()
	store := session.NewCookieStore([]byte("test"))
	store.MaxAge(60 * 60 * 12)
	e.Use(session.Sessions("pool-SSID", store))
	e.Use(permission.Permission(e, permission.Config{
		Title:    "JEX",
		Subtitle: "Permission",
		Path: permission.AuthPath{
			Login:       "/login",
			Logout:      "/logout",
			ChangePWD:   "/changepwd",
			Register:    "/register",
			CaptchaCode: "http://ethfans.org/rucaptcha",
		},
		LocalAuth: permission.LocalLoginAuth{
			AutoPage:      true,
			OnLoginBefore: OnLoginBefore,
			OnCaptchaCode: func(ec echo.Context, code string) bool {
				fmt.Println("VerifCaptchaCode: ", code)
				return false
			},
		},
		Oauth2Auth: permission.Oauth2LoginAuth{
			Enabled:      false,
			PathCallback: "/oauth2callback",
			PathError:    "/oauth2error",
		},
		PermFilter: permission.PermFilter{
			AdminPathPrefixes: []string{"/admin"},
			UserPathPrefixes:  []string{"/me"},
		},
		Cache: memory.NewMemoryCache(time.Minute * 60),
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	g := e.Group("/admin")
	g.GET("/blocks", func(c echo.Context) error {
		return c.String(http.StatusOK, "this is admin blocks page...")
	})
	g.GET("/miner", func(c echo.Context) error {
		return c.String(http.StatusOK, "this is admin miner page...")
	})
	e.GET("/me", func(c echo.Context) error {
		return c.String(http.StatusOK, "who are me ????")
	})

	// e.GET("/Confirmed", func(c echo.Context) error {
	// 	permission.SetConfirmed("admin")
	// 	return c.String(http.StatusOK, "set admin Confirmed")
	// })

	// e.GET("/UnConfirmed", func(c echo.Context) error {
	// 	permission.SetUnConfirmed("admin")
	// 	return c.String(http.StatusOK, "set admin UnConfirmed")
	// })

	e.Logger.Fatal(e.Start(":1323"))
}
