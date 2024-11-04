package qq

import (
	"admin/hack/gcmf"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/encoding/gjson"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	QAuth_Api_AuthCode    = "https://graph.qq.com/oauth2.0/authorize"
	QAuth_Api_AccessToken = "https://graph.qq.com/oauth2.0/token"
	QAuth_Api_OpenId      = "https://graph.qq.com/oauth2.0/me"
	QAuth_Api_GetUserInfo = "https://graph.qq.com/user/get_user_info"
)

type AccessTokenRes struct {
	AccessToken      string
	ExpiresIn        string
	RefreshToken     string
	Error            string
	ErrorDescription string
}
type OpenIdRes struct {
	ClientId         string
	Openid           string
	Error            string
	ErrorDescription string
}
type UserInfoRes struct {
	Ret          int
	Msg          string
	IsLost       int
	NickName     string
	FigureurlQq1 string
	Gender       string
	GenderType   string
}

var OAuth = new(oauth)

type oauth struct {
}

// GetAuthorURL 获得 LoginUrl

func (o *oauth) GetAuthorURL(Ctx context.Context, Scope ...string) string {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Qq
	_scope := "get_user_info"
	if len(Scope) > 0 {
		_scope = Scope[0]
	}
	// 将session作为State
	state, _ := g.RequestFromCtx(Ctx).Session.Id()
	_url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&display=%s",
		QAuth_Api_AuthCode, oConfig.AppID, url.QueryEscape(oConfig.RedirectUrl), _scope, state, "mobile")

	return _url

}

// GetQcURL 获取PC端连接
func (o *oauth) GetQcURL(Ctx context.Context, Scope ...string) string {
	_scope := "get_user_info"
	if len(Scope) > 0 {
		_scope = Scope[0]
	}
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Qq
	// 将session作为State
	state, _ := g.RequestFromCtx(Ctx).Session.Id()
	_url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		QAuth_Api_AuthCode, oConfig.AppID, url.QueryEscape(oConfig.RedirectUrl), _scope, state)

	return _url

}

// GetAccessToken 获得AccessToken
func (o *oauth) GetAccessToken(Ctx context.Context, authCode string) (string, error) {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Qq
	value := g.Map{
		"grant_type":    "authorization_code",
		"client_id":     oConfig.AppID,
		"redirect_uri":  oConfig.RedirectUrl,
		"client_secret": oConfig.AppKey,
		"code":          authCode,
		"fmt":           "json",
	}

	r, err := g.Client().Get(Ctx, QAuth_Api_AccessToken, value)
	if err != nil {
		return "", err
	}
	if r.Response.StatusCode != 200 {
		return "", errors.New("http code error")
	}
	defer r.Close()

	res := r.ReadAllString()

	var data *AccessTokenRes
	if j, err := gjson.DecodeToJson(res); err != nil {
		return "", err
	} else {
		if err := j.Scan(&data); err != nil {
			return "", err
		}
	}

	if data.Error != "" {

		return "", errors.New(data.Error + ":" + data.ErrorDescription)
	}

	return data.AccessToken, nil
}

// GetOpenId GetOpenId
func (o *oauth) GetOpenId(Ctx context.Context, accessToken string) (string, error) {
	value := g.Map{
		"access_token": accessToken,
		"fmt":          "json",
	}

	r, err := g.Client().Get(Ctx, QAuth_Api_OpenId, value)
	if err != nil {
		return "", err
	}
	if r.Response.StatusCode != 200 {
		return "", errors.New("http code error")
	}
	defer r.Close()

	res := r.ReadAllString()
	var data *OpenIdRes
	if j, err := gjson.DecodeToJson(res); err != nil {
		return "", err
	} else {
		if err := j.Scan(&data); err != nil {
			return "", err
		}
	}
	if data.Error != "" {
		return "", errors.New(data.Error + ":" + data.ErrorDescription)
	}
	return data.Openid, nil
}

// TokenAndOpenId a
func (o *oauth) TokenAndOpenId(Ctx context.Context, authCode string) (string, string, error) {
	accessToken, err := o.GetAccessToken(Ctx, authCode)
	if err != nil {
		return "", "", err
	}
	openId, err := o.GetOpenId(Ctx, accessToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, openId, nil
}

// GetUserInfo getUserInfo
func (o *oauth) GetUserInfo(Ctx context.Context, accessToken string, openId string) (userInfo *UserInfoRes, err error) {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Qq
	value := g.Map{
		"access_token":       accessToken,
		"oauth_consumer_key": oConfig.AppID,
		"openid":             openId,
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

	if j, err := gjson.DecodeToJson(res); err != nil {
		return nil, err
	} else {
		if err := j.Scan(&userInfo); err != nil {
			return nil, err
		}
	}
	if userInfo.Ret != 0 {
		err = errors.New(userInfo.Msg)
		return
	}
	return
}
