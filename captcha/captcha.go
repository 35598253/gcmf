package captcha

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/mojocn/base64Captcha"
)

// Store 定义储存
var Store = base64Captcha.DefaultMemStore

// KeyLong 字母长度
var KeyLong = 6

// ImgWidth 图像长度
var ImgWidth = 150

// ImgHeight 图像长度
var ImgHeight = 38

// Source 输出形式
var Source = "1234567890qwertyuioplkjhgfdsazxcvbnm"

// GetCaptcha 验证码生成
func GetCaptcha(r *ghttp.Request) (string, error) {

	driver := base64Captcha.NewDriverDigit(ImgHeight, ImgWidth, KeyLong, 0.7, 80) // 字符,公式,验证码配置, 生成默认数字的driver
	cp := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, err := cp.Generate()
	// 将ID 写入Session
	_ = r.Session.Set("CapId", id)
	return b64s, err
}

// VerifyCaptcha 验证
func VerifyCaptcha(id, answer string, clear bool) (match bool) {
	match = Store.Get(id, clear) == answer
	return
}
