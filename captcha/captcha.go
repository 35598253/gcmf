package captcha

import (
	"image/color"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/mojocn/base64Captcha"
)

const (
	keyLong   = 5
	imgWidth  = 150
	imgHeight = 38
	source    = "123456789QWERTYUIPASDFGHJKLZXCVBNM"
	captMark  = "captcha"
	noise     = 60
)

var bagColor = &color.RGBA{R: 255, G: 255, B: 255, A: 255}

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// GetMath 验证码生成
func GetMath(session *ghttp.Session, Mark ...string) (string, error) {
	var driver base64Captcha.Driver
	_cap := captMark
	if len(Mark) > 0 {
		_cap = Mark[0]
	}
	captchaConfig := base64Captcha.DriverMath{
		Height:          imgHeight,
		Width:           imgWidth,
		NoiseCount:      noise,
		ShowLineOptions: 1,
		BgColor:         bagColor,
		Fonts:           []string{"ApothecaryFont.ttf", "wqy-microhei.ttc"},
	}
	driver = captchaConfig.ConvertFonts()
	// 字符,公式,验证码配置, 生成默认数字的driver
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	// 将ID 写入Session
	_ = session.Set(_cap, id)
	return b64s, err
}

// GetString 验证码生成
func GetString() (string, string, error) {
	var driver base64Captcha.Driver
	captchaConfig := base64Captcha.DriverString{
		Height:          imgHeight,
		Width:           imgWidth,
		NoiseCount:      noise,
		ShowLineOptions: 1,
		Length:          keyLong,
		Source:          source,
		BgColor:         bagColor,

		Fonts: []string{"wqy-microhei.ttc"},
	}
	driver = captchaConfig.ConvertFonts()
	// 字符,公式,验证码配置, 生成默认数字的driver
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()

	return id, b64s, err
}

// Verify 验证
func Verify(answer string, codeId string) (match bool) {

	match = store.Get(codeId, true) == answer
	return
}
