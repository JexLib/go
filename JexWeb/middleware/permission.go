package middleware

import (
	"github.com/labstack/echo"
	"github.com/xyproto/permissionbolt"
)

type Options struct {
	DenyMessage string
	LoginPath   string
}

func PermissionMiddleware(perm *permissionbolt.Permissions, denyFunction echo.HandlerFunc) echo.MiddlewareFunc {

	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {

			if perm.Rejected(c.Response().Writer, c.Request()) {
				// Get and call the Permission Denied function
				//perm.DenyFunction()(c.Response().Writer, c.Request())
				// Reject the request by not calling the next handler below

				return denyFunction(c)
			}
			// Continue the chain of middleware
			return next(c)
		})
	})
}

// func Middleware(option Options) (echo.MiddlewareFunc, *permissionbolt.Permissions, error) {
// 	perm, err := permissionbolt.New()
// 	if perm == nil {
// 		return nil, nil, err
// 	}

// 	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return echo.HandlerFunc(func(c echo.Context) error {
// 			// Check if the user has the right admin/user rights

// 			if c.Path() != option.LoginPath && perm.Rejected(c.Response().Writer, c.Request()) {
// 				fmt.Println(option.DenyMessage)
// 				username := perm.UserState().Username(c.Request())

// 				if !perm.UserState().IsLoggedIn(username) {
// 					c.Request().Header.Set("Referer", c.Request().RequestURI)
// 					//c.Request().Referer()
// 					return c.Redirect(301, option.LoginPath+"?next="+c.Request().RequestURI)
// 				}

// 				// Deny the request, don't call other middleware handlers
// 				return echo.NewHTTPError(http.StatusForbidden, option.DenyMessage)
// 			}
// 			// Continue the chain of middleware
// 			return next(c)
// 		})
// 	}), perm, nil
// }
