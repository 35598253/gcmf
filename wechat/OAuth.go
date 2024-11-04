package wechat

import (
	"github.com/35598253/gcmf"

	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	QAuth_Api_QcUrl       = "https://open.weixin.qq.com/connect/qrconnect"
	QAuth_Api_AuthUrl     = "https://open.weixin.qq.com/connect/oauth2/authorize"
	QAuth_Api_AccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token"
	QAuth_Api_GetUserInfo = "https://api.weixin.qq.com/sns/userinfo"
)

type AccessTokenRes struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}
type UserInfoRes struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
	Errcode    int      `json:"errcode"`
	Errmsg     string   `json:"errmsg"`
}

var OAuth = new(oauth)

type oauth struct {
}

// GetAuthUrl 获取授权Url
func (o *oauth) GetAuthUrl(Ctx context.Context, Scope ...string) string {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Wechat
	// 生成授权链接
	_scope := "snsapi_login"
	if len(Scope) > 0 {
		_scope = Scope[0]
	}
	// 将session作为State
	state, _ := g.RequestFromCtx(Ctx).Session.Id()
	_url := fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		QAuth_Api_AuthUrl, oConfig.AppID, url.QueryEscape(oConfig.RedirectUrl), _scope, state)

	return _url

}

// GetQcUrl 获取授权Url
func (o *oauth) GetQcUrl(Ctx context.Context, Scope ...string) string {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Wechat
	// 生成授权链接
	_scope := "snsapi_login"
	if len(Scope) > 0 {
		_scope = Scope[0]
	}
	// 将session作为State
	state, _ := g.RequestFromCtx(Ctx).Session.Id()
	_url := fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		QAuth_Api_QcUrl, oConfig.AppID, url.QueryEscape(oConfig.RedirectUrl), _scope, state)

	return _url

}

// GetAccessToken 获取AccessToken
func (o *oauth) GetAccessToken(Ctx context.Context, code string) (*AccessTokenRes, error) {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Wechat
	value := g.Map{
		"appid":      oConfig.AppID,
		"secret":     oConfig.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	r, err := g.Client().Get(Ctx, QAuth_Api_AccessToken, value)
	if err != nil {
		return nil, err
	}
	if r.Response.StatusCode != 200 {
		return nil, errors.New("http code error")
	}
	defer r.Close()

	res := r.ReadAllString()

	var data *AccessTokenRes
	if j, err := gjson.DecodeToJson(res); err != nil {
		return nil, err
	} else {
		if err := j.Scan(&data); err != nil {
			return nil, err
		}
	}
	if data.Errcode != 0 {
		err = fmt.Errorf("errCode:%d,msg:%s", data.Errcode, data.Errmsg)
		return nil, err
	}

	return data, nil
}

// GetUserInfo 获取AccessToken
func (o *oauth) GetUserInfo(Ctx context.Context, accessToken, openId, lang string) (*UserInfoRes, error) {

	value := g.Map{
		"access_token": accessToken,
		"openid":       openId,
		"lang":         lang,
	}

	r, err := g.Client().Get(Ctx, QAuth_Api_GetUserInfo, value)
	if err != nil {
		return nil, err
	}
	if r.Response.StatusCode != 200 {
		return nil, errors.New("http code error")
	}
	defer r.Close()

	res := r.ReadAllString()

	var data *UserInfoRes
	if j, err := gjson.DecodeToJson(res); err != nil {
		return nil, err
	} else {
		if err := j.Scan(&data); err != nil {
			return nil, err
		}
	}
	if data.Errcode != 0 {
		err = fmt.Errorf("errCode:%d,msg:%s", data.Errcode, data.Errmsg)
		return nil, err
	}

	return data, nil
}
