package email

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/frame/g"

	"github.com/go-gomail/gomail"
)

// Email 服务
var Email = new(emailService)

type emailService struct {
}

// EmailConfig 邮件配置
type EmailConfig struct {
	// ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.mail.qq.com
	ServerHost string `v:"required#邮箱服务器地址不能为空"`
	// ServerPort 邮箱服务器端口，如腾讯企业邮箱为465
	ServerPort int `v:"required#邮箱服务器端口不能为空"`
	// FromName　发件人
	FromName string
	// FromEmail　发件人邮箱地址
	FromEmail string `v:"required#发件邮箱不能为空"`
	// FromPasswd 发件人邮箱密码（注意，这里是明文形式）
	FromPasswd string `v:"required#发件邮箱密码不能为空"`
	Template   []string
}

// SendEmail body支持html格式字符串
func (s *emailService) SendEmail(subject, body, ToEmail string, template int, CcEmail ...string) error {
	var cfg *EmailConfig
	if err := g.Cfg("email").ToStruct(&cfg); err != nil {
		return errors.New("读取配置出错")
	}
	// URL 通信网址

	con := fmt.Sprintf(cfg.Template[template], body)
	var tos, ccs []string
	m := gomail.NewMessage()
	for _, tmp := range strings.Split(ToEmail, ",") {
		tos = append(tos, strings.TrimSpace(tmp))
	}

	// 主题
	m.SetHeader("Subject", subject)
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
	m.SetAddressHeader("From", cfg.FromEmail, cfg.FromName)
	// 正文
	m.SetBody("text/html", con)

	d := gomail.NewDialer(cfg.ServerHost, cfg.ServerPort, cfg.FromEmail, cfg.FromPasswd)
	// 发送
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
