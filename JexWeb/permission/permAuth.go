package permission

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/JexLib/golang/cache"
	"github.com/JexLib/golang/cache/memory"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

type PermDenyFunction func(echo.Context) error       //权限不足的处理
type PermSkipFunction func(echo.Context, *User) bool //是否权限管理处理通过

type (
	//验证session用户
	OnUserAuthEvent func(user interface{}) bool
	//验证登录用户（非内置ui不需要），返回用户结构指针
	OnLoginEvent func(account, password string, us *User) error
	//注册新用户
	OnRegisterEvent func(account, password string) error
	//修改密码
	OnChangePasswordEvent func(account, password, newpassword string) error
	//验证码验证
	OnCaptchaCodeEvent func(echo.Context, string) bool
)

type Config struct {
	Title      string
	Subtitle   string
	Path       AuthPath
	LocalAuth  LocalLoginAuth  //本地登录验证
	Oauth2Auth Oauth2LoginAuth //Oauth2登录验证
	PermFilter PermFilter      //权限过滤条件
	Cookie     CookieConfig
	Cache      cache.Cache //加速并发缓存，默认使用内存缓存
}

type CookieConfig struct {
	MaxAge int    //seconds
	Domain string //cookie Domain
}

type AuthPath struct {
	Login       string
	Logout      string
	ChangePWD   string
	Register    string
	CaptchaCode string //验证码获取地址
}

//本地登录验证
type LocalLoginAuth struct {
	AutoPage         bool                  //使用本中间件自动创建页面
	OnLoginBefore    OnLoginEvent          //登录前回调
	OnRegister       OnRegisterEvent       //注册处理
	OnChangePassword OnChangePasswordEvent //修改密码处理
	OnCaptchaCode    OnCaptchaCodeEvent    //验证码检查
}

//Verif

//Oauth2登录验证
type Oauth2LoginAuth struct {
	PathCallback string
	PathError    string
	Clients      map[string]Oauth2Config
}

type Oauth2Config struct {
	oauth2.Config //支持oauth2服务器设置
	ImgUrl        string
	//Endpoints oauth2.Config
}

//权限过滤
type PermFilter struct {
	AdminPathPrefixes   []string         //必须管理员权限访问的路由前缀
	UserPathPrefixes    []string         //必须登录后访问的路由前缀
	ExcludePathPrefixes []string         //排除控制的路由前缀
	PermSkip            PermSkipFunction //自定义权限过滤函数,可以进行角色等业务的处理
	PermDeny            PermDenyFunction //权限不足时的处理函数
}

var (
	_config *Config
)

// type PermissionAuth struct {
// 	conf Config
// }

func defConfig() Config {
	return Config{
		Title:    "JEX",
		Subtitle: "Permission",
		Path: AuthPath{
			Login:     "/login",
			Logout:    "/logout",
			ChangePWD: "/changepwd",
			Register:  "/register",
		},
		LocalAuth: LocalLoginAuth{
			AutoPage: true,
		},
		Oauth2Auth: Oauth2LoginAuth{
			PathCallback: "/oauth2callback",
			PathError:    "/oauth2error",
		},
		Cache: memory.NewMemoryCache(time.Minute * 60),
	}
}

// func NewPermissionAuth(config ...Config) *PermissionAuth {
// 	if len(config) == 0 {
// 		config = append(config, defConfig())
// 	}
// 	return &PermissionAuth{
// 		conf: config[0],
// 	}
// }

func Permission(config ...Config) echo.MiddlewareFunc {
	if len(config) == 0 {
		config = append(config, defConfig())
	}
	_config = &config[0]
	newUsers(_config.Cache)
	// Oauth2AuthMiddlewares := make(map[string]echo.MiddlewareFunc)
	// if len(config[0].Oauth2Auth.Clients) > 0 {
	// 	//oauth2验证
	// 	for k, v := range config[0].Oauth2Auth.Clients {
	// 		Oauth2AuthMiddlewares[k] = oauth.NewOAuth2Provider(&v.Config)

	// 	}
	// }

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if config[0].LocalAuth.AutoPage {
				//使用内置登录界面路由
				switch c.Request().URL.Path {
				case config[0].Path.Login:
					return handle_Login(c)
				case config[0].Path.Logout:
					return handle_Logout(c)
				case config[0].Path.Register:
					return handle_Register(c)
				case config[0].Path.ChangePWD:
					return handle_ChangePwd(c)
				}
			}

			user := UserInfo(c)

			// if user.IsLogin && !user.Confirmed {
			// 	return errors.New("账号未通过审核，请与管理员联系")
			// }
			if !isExcludePath(c) {
				if isUserPath(c) && (!user.IsLogin) {
					//必须登录访问的路由，但未登录
					return c.Redirect(302, config[0].Path.Login+"?next="+base64.RawURLEncoding.EncodeToString([]byte(c.Request().RequestURI)))
				}
				if isAdminPath(c) {
					//必须admin用户访问的路由
					if !user.IsLogin {
						//必须登录
						return c.Redirect(302, config[0].Path.Login+"?next="+base64.RawURLEncoding.EncodeToString([]byte(c.Request().RequestURI)))
					}
					if !user.IsAdmin {
						//非admin，权限不够
						return handle_permDeny(c)
					}
				}
			}
			if _config.PermFilter.PermSkip != nil {
				if _config.PermFilter.PermSkip(c, user) {
					return next(c)
				} else {
					// 权限不够
					return handle_permDeny(c)
				}
			}
			return next(c)
		}
	}
}

//权限不够
func handle_permDeny(c echo.Context) error {
	if _config.PermFilter.PermDeny != nil {
		return _config.PermFilter.PermDeny(c)
	}
	return c.String(http.StatusForbidden, "Permission denied!")
}

//要排除权限控制的路径
func isExcludePath(c echo.Context) bool {
	for _, prefix := range _config.PermFilter.ExcludePathPrefixes {
		if strings.HasPrefix(c.Path(), prefix) {
			return true
		}
	}
	return false
}

func isAdminPath(c echo.Context) bool {
	for _, prefix := range _config.PermFilter.AdminPathPrefixes {
		if strings.HasPrefix(c.Path(), prefix) {
			//必须admin用户
			return true
		}
	}
	return false
}

func isUserPath(c echo.Context) bool {
	for _, prefix := range _config.PermFilter.UserPathPrefixes {
		if strings.HasPrefix(c.Path(), prefix) {
			//必须登录用户
			return true
		}
	}
	return false
}

func SetAdminPaths(paths []string) {
	_config.PermFilter.AdminPathPrefixes = paths
}

func AddAdminPath(path string) {
	_config.PermFilter.AdminPathPrefixes = append(_config.PermFilter.AdminPathPrefixes, path)
}

func SetUserPaths(paths []string) {
	_config.PermFilter.UserPathPrefixes = paths
}

func AddUserPath(path string) {
	_config.PermFilter.UserPathPrefixes = append(_config.PermFilter.UserPathPrefixes, path)
}

func handle_Login(c echo.Context) error {
	cookIfo := GetCookieInfo(c)
	if cookIfo != nil && cookIfo.IsLogin {
		return c.Redirect(302, "/")
	}
	if c.QueryParam("oauth") != "" && len(_config.Oauth2Auth.Clients) > 0 {

	}

	switch c.Request().Method {
	case "GET":
		tmpl, _ := template.New("Layout").Parse(html_layout)
		tmpl.New("SCRIPT").Parse(html_script)
		tmpl.New("STYLE").Parse(styleStr)
		tmpl.New("yeld").Parse(html_sing_in)
		data := map[string]interface{}{
			"Title":           _config.Title,
			"Subtitle":        _config.Subtitle,
			"Sign_up":         _config.Path.Register,
			"PathCaptchaCode": _config.Path.CaptchaCode,
			"Oauth2Clients":   _config.Oauth2Auth.Clients,
		}

		out := new(bytes.Buffer)
		err := tmpl.Execute(out, data)
		if err != nil {
			panic(err)
		}
		return c.HTML(200, out.String())
	case "POST":
		account := c.FormValue("user")
		passwd := c.FormValue("paswd")
		captchaCode := c.FormValue("_rucaptcha")

		if strings.Trim(account, "") == "" || strings.Trim(passwd, "") == "" || strings.Trim(captchaCode, "") == "" {
			return c.String(500, "登录失败,用户名、密码、验证码不允许为空！")
		}
		if !_config.LocalAuth.OnCaptchaCode(c, captchaCode) {
			return c.JSON(500, map[string]interface{}{
				"status":  "E-CaptchaCode",
				"message": "验证码错误，请重新输入！",
			})
		}
		user := UserInfo(c)

		if err := _config.LocalAuth.OnLoginBefore(account, passwd, user); err == nil {
			Login(c, account)
			return c.JSON(200, map[string]interface{}{
				"status":   302,
				"message":  "登录成功",
				"location": getNext(c),
			})
		} else {
			//错误消息
			return c.String(500, err.Error())
		}

	}

	return nil
}

func handle_Logout(c echo.Context) error {

	sing_out_str := html_script + html_logout
	switch c.Request().Method {
	case "GET":
		return c.HTML(200, sing_out_str)
	case "POST":
		Logout(c)
		return c.JSON(200, map[string]interface{}{
			"status":   302,
			"message":  "注销用户登录完成",
			"location": "/",
		})
	}
	return nil
}

func handle_Register(c echo.Context) error {
	switch c.Request().Method {
	case "GET":
		tmpl, _ := template.New("Layout").Parse(html_layout)
		tmpl.New("SCRIPT").Parse(html_script)
		tmpl.New("STYLE").Parse(styleStr)
		tmpl.New("yeld").Parse(html_sing_up)
		data := map[string]interface{}{
			"Title":           _config.Title,
			"Subtitle":        _config.Subtitle,
			"Sign_in":         _config.Path.Login,
			"PathCaptchaCode": _config.Path.CaptchaCode,
		}

		out := new(bytes.Buffer)
		err := tmpl.Execute(out, data)
		if err != nil {
			panic(err)
		}
		return c.HTML(200, out.String())
	case "POST":
		c.Request().ParseForm()
		account := c.FormValue("user")
		passwd := c.FormValue("paswd")
		passwd1 := c.FormValue("paswd1")
		captchaCode := c.FormValue("_rucaptcha")
		if strings.Trim(account, "") == "" || strings.Trim(passwd, "") == "" || strings.Trim(passwd1, "") == "" {
			return c.String(500, "用户名、密码不允许为空！")
		}
		if passwd != passwd1 {
			return c.String(500, "两次密码不匹配！")
		}

		if len(passwd) < 8 {
			return c.String(500, "密码必须是8位长度的数字或字符串！")
		}

		if !_config.LocalAuth.OnCaptchaCode(c, captchaCode) {
			return c.JSON(500, map[string]interface{}{
				"status":  "E-CaptchaCode",
				"message": "验证码错误，请重新输入！",
			})
		}
		if err := _config.LocalAuth.OnRegister(account, passwd); err == nil {
			Logout(c)
			return c.JSON(200, map[string]interface{}{
				"status":   302,
				"message":  "注册新用户成功! 现在将跳转到登录页面...",
				"location": _config.Path.Login,
			})
		} else {
			//错误消息
			return c.String(500, err.Error())
		}
	}
	return nil
}

//修改密码
func handle_ChangePwd(c echo.Context) error {

	switch c.Request().Method {
	case "GET":
		tmpl, _ := template.New("Layout").Parse(html_layout)
		tmpl.New("SCRIPT").Parse(html_script)
		tmpl.New("STYLE").Parse(styleStr)
		tmpl.New("yeld").Parse(html_changepwd)
		data := map[string]interface{}{
			"Title":           _config.Title,
			"Subtitle":        _config.Subtitle,
			"Sign_in":         _config.Path.Login,
			"PathCaptchaCode": _config.Path.CaptchaCode,
		}

		out := new(bytes.Buffer)
		err := tmpl.Execute(out, data)
		if err != nil {
			panic(err)
		}
		return c.HTML(200, out.String())

	case "POST":
		c.Request().ParseForm()
		account := c.FormValue("user")
		passwd := c.FormValue("paswd")
		newpasswd := c.FormValue("newpaswd")
		newpasswd1 := c.FormValue("newpaswd1")
		if strings.Trim(account, "") == "" || strings.Trim(passwd, "") == "" || strings.Trim(newpasswd, "") == "" {
			return c.String(500, "用户名、密码不允许为空！")
		}
		if newpasswd != newpasswd1 {
			return c.String(500, "两次新密码不匹配！")
		}

		if len(newpasswd) < 8 {
			return c.String(500, "密码必须是8位长度的数字或字符串！")
		}

		if err := _config.LocalAuth.OnChangePassword(account, passwd, newpasswd); err == nil {

			return c.JSON(200, map[string]interface{}{
				"status":   302,
				"message":  "修改密码成功!",
				"location": _config.Path.Login,
			})
		} else {
			//错误消息
			return c.String(500, err.Error())
		}
	}
	return nil
}

// func AuthPermission(e *echo.Echo, config ...Config) echo.MiddlewareFunc {
// 	return nil
// }

//成功登录,写session
func Login(c echo.Context, account string) {
	// sess := session.Default(c)
	// sess.Options(session.Options{
	// 	Path:     "",
	// 	Domain:   _config.Cookie.Domain,
	// 	MaxAge:   _config.Cookie.MaxAge,
	// 	HttpOnly: true,
	// })
	// sess.Save()
	cookIfo := GetCookieInfo(c)
	if cookIfo == nil {
		cookIfo = &cookieInfo{
			Account: account,
		}
	}

	cookIfo.LastLoggedTime = time.Now().Unix()
	cookIfo.IsLogin = true
	cookIfo.Save(c)
}

//退出登录,写session
func Logout(c echo.Context) {
	cookIfo := GetCookieInfo(c)
	if cookIfo == nil {
		return
	}
	cookIfo.IsLogin = false
	cookIfo.Save(c)
}

//复位当前用户,写session
// func RestartUser(c echo.Context) {

// 	//删除缓存信息，用户下次进入重新登录
// 	GetUser(c).Delete(c)
// }

//获取登录前页面地址
func getNext(c echo.Context) string {
	if next, err := base64.RawURLEncoding.DecodeString(c.QueryParam("next")); err == nil {
		n := string(next)
		if n == "" || n == _config.Path.Login {
			n = "/"
		}
		return n
	}
	return "/"
}
