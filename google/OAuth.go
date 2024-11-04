package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/35598253/gcmf"
	"net/url"

	"github.com/gogf/gf/v2/util/gconv"

	"google.golang.org/api/idtoken"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

var OAuth = new(oauth)

type oauth struct {
}
type AccessTokenRes struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	Error            string
	ErrorDescription string
}
type UserInfo struct {
	NickName      string `json:"nick_name"`
	Gid           string `json:"gid"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Picture       string `json:"picture"`
}

const (
	qAuth_Api_AuthUrl = "https://accounts.google.com/o/oauth2/v2/auth"
	qAuth_Api_Token   = "https://oauth2.googleapis.com/token"
	scope             = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
)

// GetAuthorURL 获得 LoginUrl
func (o *oauth) GetAuthorURL(Ctx context.Context, Scope ...string) string {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Google
	_scope := url.QueryEscape(scope)
	// 将session作为State
	state, _ := g.RequestFromCtx(Ctx).Session.Id()
	_url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		qAuth_Api_AuthUrl, oConfig.ClientID, url.QueryEscape(oConfig.RedirectURL), _scope, state)
	return _url
}

// GetInfo 获得GetInfo
func (o *oauth) GetInfo(Ctx context.Context, authCode string) (*UserInfo, error) {
	gCfg, _ := gcmf.GetConfig(Ctx)
	oConfig := gCfg.Oauth.Google
	_authCode, _ := url.QueryUnescape(authCode)
	value := g.Map{
		"client_id":     oConfig.ClientID,
		"client_secret": oConfig.ClientSecret,
		"redirect_uri":  oConfig.RedirectURL,
		"code":          _authCode,
		"grant_type":    "authorization_code",
	}

	r, err := g.Client().Post(Ctx, qAuth_Api_Token, value)
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

	if data.Error != "" {
		return nil, errors.New(data.Error + ":" + data.ErrorDescription)
	}

	_tID, _ := idtoken.Validate(Ctx, data.IDToken, oConfig.ClientID)
	var info *UserInfo
	_ = gconv.Struct(_tID, &info)
	return info, nil
}
