package uauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JexLib/golang/JexWeb/captcha"
	"github.com/JexLib/golang/JexWeb/session"
	"github.com/labstack/echo"
)

type (
	//验证session用户
	OnUserAuthEvent func(user interface{}) bool
	//验证登录用户（非内置ui不需要），返回用户结构指针
	OnLoginEvent func(account, password string) (interface{}, error)
	//注册新用户
	OnRegisterEvent func(account, password string) error
	//修改密码
	OnChangePasswordEvent func(account, password, newpassword string) error

	Config struct {
		//Skipper      middleware.Skipper
		UseBuiltinUi bool //是否使用内置UI
		Title        string
		LoginTitle   string
		MaxAge       int    //seconds
		Domain       string //cookie Domain
		Event        Events
		Path         Paths
	}

	Paths struct {
		SignIN    string
		SignOUT   string
		SignUP    string
		ChangePWD string
	}
	Events struct {
		OnUserAuth       OnUserAuthEvent
		OnLogin          OnLoginEvent
		OnRegister       OnRegisterEvent
		OnChangePassword OnChangePasswordEvent
	}

	// User interface{
	// 	IsAdmin() bool
	// 	IsLogin() bool

	// }
	// UAUser struct {
	// 	LoginTime int64
	// 	IsLogin  bool
	// 	Admin     bool        //是否管理员
	// 	Roles     interface{} //角色信息
	// 	Extend    interface{} //其他信息
	// }

	UAuthInfo struct {
		LoginTime int64
		IsLogin   bool
		User      interface{}
	}
)

var (
	_config            *Config
	_key               = "auth_info"
	_CaptchaCodeServer *captcha.CaptchaServer
)

// type testuser struct {
// 	id   int
// 	name string
// }

func defConfing() Config {
	return Config{
		UseBuiltinUi: true,
		Title:        "JEX",
		LoginTitle:   "登录系统",
		MaxAge:       60 * 60 * 12, //默认12小时
		Path: Paths{
			SignIN:    "/login",
			SignOUT:   "/logout",
			SignUP:    "/register",
			ChangePWD: "/changepwd",
		},
	}
}

func UAuthMiddleware(e *echo.Echo, config ...Config) echo.MiddlewareFunc {
	if len(config) == 0 {
		config = append(config, defConfing())
	}
	if config[0].Event.OnUserAuth == nil || config[0].Event.OnLogin == nil || config[0].Event.OnRegister == nil || config[0].Event.OnChangePassword == nil {
		fmt.Println("must have OnUserAuthEvent and OnLoginEvent,OnRegister,OnChangePassword func. ")
		os.Exit(0)
	}
	//_CaptchaCodeServer = captcha.NewCaptchaServer(e, "/captcha", 4, "public/assets/fonts/comic.ttf", "public/assets/fonts/D3Parallelism.ttf")
	_config = &config[0]
	mfunc := func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			if !verifySession(config[0], c) {
				// if c.Path() != config[0].SignIN && c.Path() != config[0].SignOUT && c.Path() != config[0].SignUP {

				if c.Path() != config[0].Path.SignIN {
					return c.Redirect(302, config[0].Path.SignIN+"?next="+base64.RawURLEncoding.EncodeToString([]byte(c.Request().RequestURI)))
				}
				// }
			} else {
				if c.Path() == config[0].Path.SignIN {
					return c.Redirect(302, getNext(c))
				}
			}
			return next(c)
		}
	}

	if config[0].UseBuiltinUi {
		//注册内置登录界面路由
		e.Any(config[0].Path.SignIN, handle_Login, mfunc)
		e.Any(config[0].Path.SignOUT, handle_Logout)
		e.Any(config[0].Path.SignUP, handle_Register)
		e.Any(config[0].Path.ChangePWD, handle_ChangePwd)

	}

	return mfunc
}

//验证session信息
func verifySession(cnf Config, c echo.Context) bool {
	sess := session.Default(c)
	s := sess.Get(_key)
	if s == nil {
		return false
	}
	if str, ok := sess.Get(_key).(string); ok {
		if storeinfo, err := storeinfoDecode(str); err == nil {
			if cnf.Event.OnUserAuth != nil {
				if cnf.Event.OnUserAuth(storeinfo.User) {
					storeinfo.IsLogin = true
					sess.Save()
					return true
				}
				//验证失败，删除session记录
				sess.Delete(_key)
				sess.Save()
				return false
			}
		}
	}
	return false
}

//成功登录,写session
func Login(user interface{}, c echo.Context) {
	sess := session.Default(c)
	sess.Options(session.Options{
		Path:     "",
		Domain:   _config.Domain,
		MaxAge:   _config.MaxAge,
		HttpOnly: true,
	})
	uAuthInfo := UAuthInfo{
		LoginTime: time.Now().Unix(),
		IsLogin:   true,
		User:      user,
	}

	//sess..Options.MaxAge = _config.MaxAge
	sess.Set(_key, uAuthInfo.encode())
	sess.Save()
}

//退出登录，删除session
func Logout(c echo.Context) {
	sess := session.Default(c)
	sess.Delete(_key)
	sess.Save()
}

func GetUAuthInfo(c echo.Context) *UAuthInfo {
	sess := session.Default(c)
	if str, ok := sess.Get(_key).(string); ok {
		if storeinfo, err := storeinfoDecode(str); err == nil {
			return storeinfo
		}
	}

	return nil
}

func handle_Login(c echo.Context) error {
	sing_in_str := strings.Replace(html_layout, "{{.yeld}}", html_sing_in, 1)
	sing_in_str = strings.Replace(sing_in_str, "{{.sign_up}}", _config.Path.SignUP, 1)
	sing_in_str = strings.Replace(sing_in_str, "{{.logintitle1}}", _config.Title, 1)
	sing_in_str = strings.Replace(sing_in_str, "{{.logintitle2}}", _config.LoginTitle, 1)
	sing_in_str = styleStr + html_script + sing_in_str
	fmt.Println(sing_in_str)
	switch c.Request().Method {
	case "GET":
		return c.HTML(200, sing_in_str)
	case "POST":
		account := c.FormValue("user")
		passwd := c.FormValue("paswd")
		captchaCode := c.FormValue("_rucaptcha")
		if strings.Trim(account, "") == "" || strings.Trim(passwd, "") == "" {
			return c.String(500, "登录失败,用户名、密码不允许为空！")
		}

		if !_CaptchaCodeServer.Validation(c, captchaCode) {
			return c.String(500, "验证码错误，请重新输入！")
		}
		if user, err := _config.Event.OnLogin(account, passwd); err == nil {
			Login(user, c)

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

//获取登录前页面地址
func getNext(c echo.Context) string {
	if next, err := base64.RawURLEncoding.DecodeString(c.QueryParam("next")); err == nil {
		n := string(next)
		if n == "" || n == _config.Path.SignIN {
			n = "/"
		}
		return n
	}
	return "/"
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
	sing_up_str := strings.Replace(html_layout, "{{.yeld}}", html_sing_up, 1)
	sing_up_str = strings.Replace(sing_up_str, "{{.sign_in}}", _config.Path.SignIN, 1)
	sing_up_str = strings.Replace(sing_up_str, "{{.logintitle1}}", _config.Title, 1)
	sing_up_str = strings.Replace(sing_up_str, "{{.logintitle2}}", _config.LoginTitle, 1)
	sing_up_str = styleStr + html_script + sing_up_str

	switch c.Request().Method {
	case "GET":
		return c.HTML(200, sing_up_str)
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
		if !_CaptchaCodeServer.Validation(c, captchaCode) {
			return c.String(500, "验证码错误，请重新输入！")
		}
		if err := _config.Event.OnRegister(account, passwd); err == nil {
			Logout(c)
			return c.JSON(200, map[string]interface{}{
				"status":   302,
				"message":  "注册新用户成功! 现在将跳转到登录页面...",
				"location": _config.Path.SignIN,
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
	change_str := strings.Replace(html_layout, "{{.yeld}}", html_changepwd, 1)
	change_str = strings.Replace(change_str, "{{.sign_in}}", _config.Path.SignIN, 1)
	change_str = strings.Replace(change_str, "{{.logintitle1}}", _config.Title, 1)
	change_str = strings.Replace(change_str, "{{.logintitle2}}", _config.LoginTitle, 1)
	change_str = styleStr + html_script + change_str

	switch c.Request().Method {
	case "GET":
		return c.HTML(200, change_str)
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

		if err := _config.Event.OnChangePassword(account, passwd, newpasswd); err == nil {
			return c.JSON(200, map[string]interface{}{
				"status":   302,
				"message":  "修改密码成功!",
				"location": _config.Path.SignIN,
			})
		} else {
			//错误消息
			return c.String(500, err.Error())
		}
	}
	return nil
}

func (st *UAuthInfo) encode() string {
	bytes, _ := json.Marshal(st)
	return base64.StdEncoding.EncodeToString(bytes)
}

func storeinfoDecode(str string) (*UAuthInfo, error) {
	if bytes, err := base64.StdEncoding.DecodeString(str); err == nil {
		st := &UAuthInfo{}
		json.Unmarshal(bytes, st)
		return st, nil
	} else {
		return nil, err
	}
}
