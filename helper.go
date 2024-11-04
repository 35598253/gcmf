package gcmf

import (
	"admin/hack/gcmf/email"
	"admin/hack/gcmf/sign"
	"admin/hack/gcmf/ztcms"
	"unicode/utf8"

	"github.com/gogf/gf/v2/i18n/gi18n"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/encoding/gbase64"

	"github.com/gogf/gf/v2/crypto/gmd5"

	"github.com/gogf/gf/v2/os/gtime"

	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

// GetCode 生成session 随机Code
func GetCode(Ctx context.Context, key string, duration time.Duration) (int, error) {
	ok, _ := gcache.Contains(Ctx, key)
	if ok {
		return 8888, errors.New("Code_Exist")
	}
	code := grand.N(10000, 999999)
	_ = gcache.Set(gctx.New(), key, code, duration)
	return code, nil
}

// VerifyCode 验证 Code
func VerifyCode(Ctx context.Context, Key, Code string) bool {

	ok, _ := gcache.Contains(Ctx, Key)
	if ok {
		c, _ := gcache.Get(Ctx, Key)

		if c.String() == Code {
			if _, e := gcache.Remove(Ctx, Key); e == nil {
				return true
			}

		}
	}
	return false
}

// IsAjaxPost AjAX提交
func IsAjaxPost(r *ghttp.Request) bool {

	return r.Request.Method == "POST" && r.IsAjaxRequest()
}

// IsPost Post提交
func IsPost(r *ghttp.Request) bool {

	return r.Request.Method == "POST"
}

// IsAjaxData layData 数据获取
func IsAjaxData(r *ghttp.Request) bool {
	data := r.GetHeader("layData")
	return data == "true"
}

// JSONToSlice JSON转化
func JSONToSlice(JSON string) []string {
	s := []byte(JSON)
	var w []interface{}
	_ = json.Unmarshal(s, &w)
	return gconv.SliceStr(w)
}

// DbCache // 粗暴清楚数据缓存
func DbCache(Dbname string) error {
	// 粗暴清除数据缓存
	if err := g.DB(Dbname).GetCache().Clear(gctx.New()); err != nil {
		return err
	}
	return nil
}

// JSONToMap JSON转化
func JSONToMap(JSON string) map[string]interface{} {
	s := []byte(JSON)
	var w []interface{}
	_ = json.Unmarshal(s, &w)
	return gconv.Map(w)
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

// InArray 切片检测
func InArray(Value string, Array []string) bool {
	for _, v := range Array {
		if Value == v {
			return true
		}
	}
	return false
}

// MergeMapStr Map合并
func MergeMapStr(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// CheckMobile 检测手机号
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}
func CheckEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// GetPass 获取密码
func GetPass(data string, Cost ...int) string {
	_cost := 10
	if len(Cost) > 0 {
		_cost = Cost[0]
	}
	pass2, _ := bcrypt.GenerateFromPassword([]byte(data), _cost)
	return gconv.String(pass2)
}

// CleanCacheAll 清理全部缓存
func CleanCacheAll(Ctx context.Context) error {
	keys, _ := gcache.Keys(Ctx)
	return gcache.Removes(Ctx, keys)
}

// GetSn 订单号码
func GetSn(Prefix ...string) string {
	var _pre string
	if len(Prefix) > 0 {
		_pre = Prefix[0]
	}
	t := gtime.Now()
	return _pre + t.TimestampMicroStr()
}

//// GetApiToken 根据时间获取Token
//func GetApiToken(Uid uint) string {
//	// 先进行md5 加密
//
//	t := gtime.Now()
//	_ut := gconv.String(Uid) + "-" + t.TimestampMicroStr()
//	Token, _ := gmd5.Encrypt(_ut)
//	return Token
//}

// GetToken 获取Token
func GetToken(Uid uint, Key string) string {
	// 先进行md5 加密
	keyMd5, _ := gmd5.Encrypt(Key)
	// 得到Token
	key := []byte(keyMd5)
	t := gtime.Timestamp()
	_ut := fmt.Sprintf("%d-%d", Uid, t)
	res, _ := gaes.Encrypt(gconv.Bytes(_ut), key)

	Token := gbase64.EncodeToString(res)
	return Token
}

// TokenGetUid 获取Uid
func TokenGetUid(Token, Key string) (uint, int64, error) {
	// 先进行md5 加密
	keyMd5, _ := gmd5.Encrypt(Key)
	// 得到Token
	key := []byte(keyMd5)
	// 处理Token
	_token, _ := gbase64.DecodeString(Token)

	_ut, err := gaes.Decrypt(_token, key)
	if err != nil {
		return 0, 0, err
	}
	arr := strings.Split(gconv.String(_ut), "-")
	return gconv.Uint(arr[0]), gconv.Int64(arr[1]), nil
}

// MapASCII MapASCII排序

func GetSignMd5(Maps interface{}, AppKey string) (string, string) {

	out, md5 := sign.GetSignMd5(Maps, "key", AppKey, "%.2f")
	//fmt.Println(_t)
	return out, md5
}

// ParseBank 隐藏银行卡中间4位
func ParseBank(Str string) string {
	let := len(Str) - 4
	reg := "\\d{" + strconv.Itoa(let) + "}(\\d{4})"
	re := regexp.MustCompile(reg)
	var str string
	for i := 0; i < len(Str)-4; i++ {
		str += "*"
	}
	str = str + "$1"
	return re.ReplaceAllStringFunc(Str, func(m string) string {
		return re.ReplaceAllString(m, str)
	})

}

// ParsePhone 隐藏手机号中间4位
func ParsePhone(Phone string) string {

	return Phone[:3] + "****" + Phone[7:]

}

// ParseEmail 隐藏邮箱中间4位
func ParseEmail(Email string) string {

	atIndex := strings.Index(Email, "@")
	if atIndex == -1 || atIndex <= 1 {
		return Email
	}

	username := Email[:atIndex]
	domain := Email[atIndex:]

	var builder strings.Builder
	builder.WriteByte(username[0])
	builder.WriteString(strings.Repeat("*", len(username)-2))
	builder.WriteByte(username[len(username)-1])

	return builder.String() + domain

}

// FormatFloat 四舍五入，ROUND_HALF_UP 模式实现
func FormatFloat(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}

// DateSubDay 计算天数

func DateSubDay(Start, End *gtime.Time) float64 {

	day := End.Sub(Start).Hours() / 24

	return day
}

// FirstStr 获取第一个字符
func FirstStr(Str string) string {

	r, _ := utf8.DecodeRuneInString(Str)
	return string(r)
}

// GetLowStr 获取小写字符串并去除空格
func GetLowStr(Str string) string {
	_str := strings.TrimSpace(Str)
	return strings.ToLower(_str)
}

// RestartServer 重启
func RestartServer(Ctx context.Context, AppPath ...string) (err error) {
	path := os.Args[0]

	if len(AppPath) > 0 {
		path = AppPath[0]
	}
	// IDE 中不重启
	if !strings.Contains(path, "/tmp/") {
		err = ghttp.RestartAllServer(Ctx, path)
	}
	return
}

// TableExist 数据表是否存在
func TableExist(Ctx context.Context, Table string, Site ...string) bool {
	Db := g.DB()
	if len(Site) > 0 {
		Db = g.DB(Site[0])
	}
	sql := "SHOW TABLES LIKE '" + Table + "'"

	res, _ := Db.Query(Ctx, sql)

	return !res.IsEmpty()

}

// TableFieldExist 数据表是否存在
func TableFieldExist(Ctx context.Context, TableName, FieldName string, Site ...string) bool {
	Db := g.DB()
	if len(Site) > 0 {
		Db = g.DB(Site[0])
	}
	sql := "SHOW COLUMNS FROM `" + TableName + "` LIKE " + "'" + FieldName + "'"

	res, _ := Db.Query(Ctx, sql)

	return !res.IsEmpty()

}

// ErrLang 统一全站多语言
func ErrLang(Str string, values ...interface{}) error {

	if len(values) > 0 {
		sArr := gconv.SliceStr(values)
		return errors.New(Str + "," + strings.Join(sArr, ","))
	} else {
		return errors.New(Str)
	}
}

// GetLang 获取多语言
func GetLang(Ctx context.Context, Message string) string {
	_msg := strings.Split(Message, ",")

	if len(_msg) == 1 {
		return gi18n.T(Ctx, _msg[0])
	}
	return gi18n.Tf(Ctx, _msg[0], _msg[1:])
}

// SendSms 发送信息
func SendSms(Ctx context.Context, Phone, Content string) error {
	// 检测是否发送
	gConfig, err := GetConfig(Ctx)

	if err != nil {
		return errors.New("Config_Is_Err")
	}
	_cfg := gConfig.Sms

	sms := ztcms.New(_cfg.Config)
	// 发送短信
	if err := sms.SendSms(Ctx, Phone, Content); err != nil {
		return err
	}
	return nil
}

// SendEmail 发送邮件
func SendEmail(Ctx context.Context, Email, Title, Body string) error {

	gCfg, _ := GetConfig(Ctx)

	_email := email.New(gCfg.Email.Config)
	tt := _email.SendEmail(Email, Title, Body)
	return tt
}
