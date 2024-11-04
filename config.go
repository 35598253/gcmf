package gcmf

import (
	"admin/hack/gcmf/email"
	"admin/hack/gcmf/ztcms"
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Upload UploadCfg
	Ip     IpConfig
	Email  EmailConfig
	Oauth  OauthConfig
	Sms    SmsConfig
}
type EmailConfig struct {
	Config   *email.Config
	Template *ETemplate
}
type ETemplate struct {
	Yzm string
}
type IpConfig struct {
	Type    int
	AmapKey string
	Ip2Path string
}
type SmsConfig struct {
	Config   *ztcms.Config
	TempLate *SmsTempLate
}
type SmsTempLate struct {
	Yzm string `json:"zym"`
}

type OauthConfig struct {
	Wechat *OauthCC
	Qq     *QqCC
	Google *GoogleCC
}
type OauthCC struct {
	AppID       string
	AppSecret   string
	RedirectUrl string
	Status      bool
}
type QqCC struct {
	AppID       string
	AppKey      string
	RedirectUrl string
	Status      bool
}
type GoogleCC struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Status       bool
}

func GetConfig(Ctx context.Context) (config Config, err error) {

	t := gfile.IsFile("./config/gcmf.toml")
	_upload := UploadCfg{
		FileExts:    "jpg|jpeg|gif|png|zip|7z|doc|docx|mp3|mp4|avi|mpg|mov|rm|rmvb|wps|txt|xlsx|xls|ppt|pptx",
		ImgExts:     "jpg|jpeg|gif|png",
		MaxSize:     8 * 1024 * 1024,
		ImageResize: 1920,
	}

	//_email := EmailConfig{
	//	FromName:   "",
	//	FromEmail:  "",
	//	FromPasswd: "",
	//	ServerHost: "",
	//	Port:       465,
	//	Template:   &ETempLate{Yzm: ""},
	//}
	_oauth := OauthConfig{
		Wechat: &OauthCC{
			AppID:       "",
			AppSecret:   "",
			RedirectUrl: "",
		},
		Qq: &QqCC{
			AppID:       "",
			AppKey:      "",
			RedirectUrl: "",
		},
		Google: &GoogleCC{
			ClientID:     "",
			ClientSecret: "",
			RedirectURL:  "",
		},
	}
	_sms := SmsConfig{
		TempLate: &SmsTempLate{
			Yzm: "",
		},
	}
	_email := EmailConfig{
		Config: &email.Config{
			FromName:   "",
			FromEmail:  "",
			FromPasswd: "",
			ServerHost: "",
			Port:       0,
		},
		Template: &ETemplate{Yzm: ""},
	}
	//配置文件文件不存在
	if !t {
		//写入默认配置
		buf, err := os.Create("./config/gcmf.toml")
		if err != nil {
			return config, err
		}
		config = Config{Upload: _upload, Oauth: _oauth, Sms: _sms, Email: _email}

		if err := toml.NewEncoder(buf).Encode(config); err != nil {
			return config, err
		}
		defer buf.Close()

		return config, nil
	}

	//配置存在,但选项没有
	res, _ := g.Cfg("gcmf").Data(Ctx)

	_ = gconv.Structs(res, &config)

	return
}
