package oauth

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/util/gconv"
)

// Token 成功返回格式
type Token struct {
	AccessToken  string `json:"access_token"`            // 网页授权接口调用凭证
	ExpiresIn    int64  `json:"expires_in"`              // access_token 接口调用凭证超时时间, 单位: 秒
	RefreshToken string `json:"refresh_token,omitempty"` // 刷新 access_token 的凭证
	OpenID       string `json:"openid,omitempty"`
	UnionID      string `json:"unionid,omitempty"`
	Scope        string `json:"scope,omitempty"` // 用户授权的作用域, 使用逗号(,)分隔
}

// UserInfo 用户信息
type UserInfo struct {
	OpenID   string `json:"openid"`   // 用户的唯一标识
	Nickname string `json:"nickname"` // 用户昵称
	Sex      int    `json:"sex"`      // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	City     string `json:"city"`     // 普通用户个人资料填写的城市
	Province string `json:"province"` // 用户个人资料填写的省份
	Country  string `json:"country"`  // 国家, 如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImageURL string `json:"headimgurl,omitempty"`

	Privilege []string `json:"privilege,omitempty"` // 用户特权信息，json 数组，如微信沃卡用户为
	UnionID   string   `json:"unionid,omitempty"`   // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
}

// MpConfig 配置
type MpConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
}

var mpCfg *MpConfig

// AuthCodeURL 生成网页授权地址.
//  appID:       公众号的唯一标识
//  redirectURI: 授权后重定向的回调链接地址
//  scope:       应用授权作用域
//  state:       重定向后会带上 state 参数, 开发者可以填写 a-zA-Z0-9 的参数值, 最多128字节
func AuthCodeURL(state ...string) string {
	_ = g.Cfg("wechat").GetStruct("mp", &mpCfg)
	sta := "login"
	if len(state) > 0 {
		sta = state[0]
	}
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + mpCfg.AppID +
		"&redirect_uri=" + mpCfg.RedirectURI +
		"&response_type=code&scope=snsapi_userinfo" +
		"&state=" + sta +
		"#wechat_redirect"
}

// GetAccessToken 换取AccessTokenToken
func GetAccessToken(Token string) (*Token, error) {
	_ = g.Cfg("wechat").GetStruct("mp", &mpCfg)
	_url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + mpCfg.AppID + "&secret=" + mpCfg.AppSecret + "&code=" + Token + "&grant_type=authorization_code"
	return getAccToken(_url)
}

// RefreshToken 刷新Token
func RefreshToken(appID, refreshToken string) (*Token, error) {
	_url := "https://open.weixin.qq.com/connect/oauth2/refresh_token?appid=" + appID + "&grant_type=" + refreshToken + "&refresh_token=REFRESH_TOKEN"

	return getAccToken(_url)
}

// GetUserInfo 获取信息
func GetUserInfo(token *Token, lang ...string) (*UserInfo, error) {
	var lan = "zh_CN"
	if len(lang) > 0 {
		lan = lang[0]
	}
	_url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + token.AccessToken + "&openid=" + token.OpenID + "&lang=" + lan
	resp, err := http.Get(_url)
	if err != nil {
		return nil, err
	}
	var userInfo *UserInfo
	body, _ := ioutil.ReadAll(resp.Body)
	bd := gconv.Map(body)
	if _, k := bd["errcode"]; k {
		return nil, errors.New(gconv.String(body))
	}
	_ = gconv.Struct(body, &userInfo)
	if userInfo == nil {
		return nil, errors.New(gconv.String(body))
	}
	return userInfo, nil
}

// getAccToken 获取token
func getAccToken(url string) (*Token, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var acToken *Token
	body, _ := ioutil.ReadAll(resp.Body)
	bd := gconv.Map(body)
	if _, k := bd["errcode"]; k {
		return nil, errors.New(gconv.String(body))
	}
	_ = gconv.Struct(body, &acToken)

	return acToken, nil
}
