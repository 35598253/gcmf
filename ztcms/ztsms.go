package ztcms

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
)

type Config struct {
	UserName string
	PassWord string
}

// config Email 配置
type config struct {
	*Config
}

func New(cfg *Config) Loader {
	return &config{
		cfg,
	}
}

type Loader interface {
	SendSms(Ctx context.Context, phone, content string) error
}

// SendSms 发送短信
func (c *config) SendSms(Ctx context.Context, phone, content string) error {

	var url = "https://hy.mix2.zthysms.com/sendSms.do"
	key, s := c.getSmsKey()

	res := g.Client().ContentType("application/x-www-form-urlencoded").PostContent(Ctx, url, g.Map{
		"username": c.UserName,
		"password": s,
		"tkey":     key,
		"mobile":   phone,
		"content":  content,
	})
	b := strings.Split(res, ",")

	if b[0] != "1" {
		return errors.New(b[1])
	}
	return nil
}

// GetSmsNum 查询剩余短信条数
func (c *config) GetSmsNum() (interface{}, error) {

	//// URL 通信网址
	//
	//var url = "https://hy.mix2.zthysms.com/balance.do"
	//key, s := getSmsKey(cfg)
	//resp, err := http.Post(url,
	//	"application/x-www-form-urlencoded",
	//	strings.NewReader("username="+cfg.Username+"&password="+s+"&tkey="+key))
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	////
	//sb := gconv.String(body)
	//b := strings.Contains(sb, ",")
	//if b {
	//	return nil, errors.New(sb)
	//}
	//
	//return sb, nil
	return nil, nil
}

// getSmsKey
func (c *config) getSmsKey() (string, string) {

	key := time.Now().Format("20060102150405")
	pass, _ := gmd5.Encrypt(c.PassWord)
	s, _ := gmd5.Encrypt(pass + key)
	return key, s
}
