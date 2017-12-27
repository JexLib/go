package captcha

import (
	"fmt"
	"image/color"
	"image/png"
	"strings"

	"github.com/JexLib/golang/JexWeb/session"
	"github.com/JexLib/golang/utils"

	"github.com/afocus/captcha"
	"github.com/labstack/echo"
)

type CaptchaServer struct {
	captcha *Captcha
	length  int
	cfg     Config
}

type Config struct {
	Path        string        //验证码URLPath
	CharCount   int           //字符数量
	ImgWidth    int           //图片宽度
	ImgHeight   int           //图片高度
	FontDir     string        //字体文件目录
	Background  []color.Color //图片背景颜色 默认白色
	CharColor   []color.Color //字符颜色数组
	CharSize    float64       //文字高度占图片高度比例，默认0.8
	DisturLevel DisturLevel   //干扰程度
	IgnoreCase  bool          //是否忽略大小写
}

//开启验证码服务
func NewCaptchaServer(e *echo.Echo, cnf Config) (*CaptchaServer, error) {
	fonts, err := utils.ListDir(cnf.FontDir, ".ttf")
	if err != nil {
		fmt.Println("read FontDir error:", err)
		return nil, err
	}
	cap := New()
	cap.SetFontSize(cnf.CharSize)
	cap.SetSize(cnf.ImgWidth, cnf.ImgHeight)
	if cnf.DisturLevel == 0 {
		cnf.DisturLevel = MEDIUM
	}
	cap.SetDisturbance(cnf.DisturLevel)
	cap.SetFrontColor(cnf.CharColor...)
	cap.SetBkgColor(cnf.Background...)
	if err := cap.SetFont(fonts...); err != nil {
		panic(err.Error())
	}
	cs := &CaptchaServer{
		captcha: cap,
		length:  cnf.CharCount,
		cfg:     cnf,
	}
	if cnf.Path == "" {
		cnf.Path = "/captcha"
	}
	e.GET(cnf.Path, cs.handle_Get)
	e.POST(cnf.Path, cs.handle_Check)
	return cs, nil
}

//获取验证码图片
func (cs *CaptchaServer) handle_Get(c echo.Context) error {
	img, str := cs.captcha.Create(cs.length, captcha.ALL)
	sess := session.Default(c)
	sess.Set("@Captcha", str)
	if err := sess.Save(); err != nil {
		fmt.Println("session  write Failure:", err)
	}
	return png.Encode(c.Response().Writer, img)
}

//url 检查验证码
func (cs *CaptchaServer) handle_Check(c echo.Context) error {
	if cs.Validation(c, c.QueryParam("code")) {
		return c.String(200, "ok")
	}

	return c.String(400, "验证码错误")
}

func (cs *CaptchaServer) Validation(c echo.Context, captchaCode string) bool {
	sess := session.Default(c)
	s := sess.Get("@Captcha")
	// fmt.Println("Captcha:", s)
	if str, ok := s.(string); ok {
		fmt.Println("CaptchaCode:", captchaCode, str)
		if cs.cfg.IgnoreCase {
			return strings.ToLower(captchaCode) == strings.ToLower(str)
		}
		return captchaCode == str
	}

	return false
}
