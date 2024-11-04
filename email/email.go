package email

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type Config struct {
	FromName   string `json:"fromName"`
	FromEmail  string `json:"fromEmail"`
	FromPasswd string `json:"fromPasswd"`
	ServerHost string `json:"serverHost"`
	Port       int    `json:"port"`
}

// config Email 配置
type config struct {
	*Config
}

type Loader interface {
	SendEmail(ToEmail, Title, Body string, CcEmail ...string) error
}

func New(Cfg *Config) Loader {

	return &config{Cfg}
}

// SendEmail body支持html格式字符串
func (c *config) SendEmail(ToEmail, Title, Body string, CcEmail ...string) error {
	// URL 通信网址

	var tos, ccs []string
	m := gomail.NewMessage()
	for _, tmp := range strings.Split(ToEmail, ",") {
		tos = append(tos, strings.TrimSpace(tmp))
	}

	// 主题
	m.SetHeader("Subject", Title)
	// 设定收件人
	m.SetHeader("To", tos...)
	//抄送列表
	if len(CcEmail) != 0 {
		for _, tmp := range strings.Split(CcEmail[0], ",") {
			ccs = append(ccs, strings.TrimSpace(tmp))
		}
		m.SetHeader("Cc", ccs...)
	}
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", c.FromEmail, c.FromName)
	// 正文
	m.SetBody("text/html", Body)

	d := gomail.NewDialer(c.ServerHost, c.Port, c.FromEmail, c.FromPasswd)
	// 发送
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
