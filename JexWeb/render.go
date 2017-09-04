package jexweb

import (
	"io"

	"github.com/labstack/echo"
	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1"
)

type RenderWrapper struct { // We need to wrap the renderer because we need a different signature for echo.
	rnd *render.Render
}

func (r *RenderWrapper) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if c.Get("Layout") != nil {
		return r.rnd.HTML(w, 0, name, data, render.HTMLOptions{c.Get("Layout").(string)})
	}

	// if data != nil {
	// 	ds, ok := data.(map[string]interface{})
	// 	if ok && ds["@LAYOUT@"] != nil {
	// 		if ds["@LAYOUT@"].(string) == "" {
	// 			return r.rnd.HTML(w, 0, name, data, render.HTMLOptions{})
	// 		} else {
	// 			return r.rnd.HTML(w, 0, name, data, render.HTMLOptions{ds["@LAYOUT@"].(string)})
	// 		}
	// 	}
	// }

	return r.rnd.HTML(w, 0, name, data) // The zero status code is overwritten by echo.
}
