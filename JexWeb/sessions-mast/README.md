# sessions
sessions

# Usage

```
func main() {
	e := echo.New()
	store := sessions.NewCookieStore([]byte("123456"))
	e.Use(sessions.Sessions("echosession", store))
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/", func(c *echo.Context) error {
		session := sessions.Default(c)
		var value int
		val := session.Get("key")
		if val == nil {
			value = 0
		} else {
			value = val.(int)
			value += 1
		}
		session.Set("key", value)
		session.Save()
		return c.JSON(200, value)
	})

	e.Run(":9991")
}

```