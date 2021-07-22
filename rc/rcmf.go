package rc

import (
	"encoding/json"
	"nmrich/rcmf/core/ztcms"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/net/ghttp"

	"github.com/gogf/gf/os/gcache"

	"github.com/gogf/gf/util/grand"
)

// GetCode 生成session 随机Code
func GetCode(key string, duration time.Duration) int {

	code := grand.N(10000, 999999)
	gcache.Set(key, code, duration)
	return code
}

// VerifyCode 验证 Code
func VerifyCode(key string, code int) bool {
	ok, _ := gcache.Contains(key)
	if ok {
		c, _ := gcache.Get(key)

		if c == code {
			if _, e := gcache.Remove(key); e == nil {
				return true
			}

		}
	}
	return false
}

// IsAjaxPost AjAX提交
func IsAjaxPost(r *ghttp.Request) bool {
	if r.Request.Method == "POST" && r.IsAjaxRequest() {
		return true
	}
	return false
}

// IsAjaxData layData 数据获取
func IsAjaxData(r *ghttp.Request) bool {
	data := r.GetHeader("layData")
	return data == "true"
}

// IsDesk Ajax加载
func IsDesk(r *ghttp.Request) bool {
	return r.Request.Header.Get("urltype") == "desk"
}

// JSONToSlice JSON转化
func JSONToSlice(JSON string) []string {
	s := []byte(JSON)
	var w []interface{}
	_ = json.Unmarshal(s, &w)
	return gconv.SliceStr(w)
}

// Stripslashes 函数Json添加的反斜杠。
func Stripslashes(str string) string {
	var dstRune []rune
	strRune := []rune(str)
	strLen := len(strRune)
	for i := 0; i < strLen; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

// SendYzm 发送信息
func SendYzm(phone string) error {
	code := GetCode(phone, 2*time.Minute)

	//	fmt.Println(match)
	if err := ztcms.SendSms(phone, gconv.String(code)); err != nil {
		return err
	}
	return nil
}
