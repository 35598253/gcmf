package gcmf

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const tgAPI = "https://api.telegram.org/bot"

var Telegram = new(telegramStruct)

type telegramStruct struct {
}

type webhookRes struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}
type webhookInfoRes struct {
	Ok     bool `json:"ok"`
	Result struct {
		Url                  string `json:"url"`
		HasCustomCertificate bool   `json:"has_custom_certificate"`
		PendingUpdateCount   int    `json:"pending_update_count"`
		MaxConnections       int    `json:"max_connections"`
		IpAddress            string `json:"ip_address"`
	} `json:"result"`
}

// SetWebhook 设置监控
func (t *telegramStruct) SetWebhook(Ctx context.Context, Token, WebHookUrl string) bool {
	uri := fmt.Sprintf("%s%s/setWebhook", tgAPI, Token)
	var webRes webhookRes
	_data := g.Map{
		"url": WebHookUrl,
	}
	res := g.Client().PostContent(Ctx, uri, _data)

	_ = gconv.Struct(res, &webRes)
	if !webRes.Ok {
		g.Log().Stdout(true).Cat("Tg").Print(Ctx, "Webhook:SetWebhook Error"+webRes.Description)
		return false
	}
	if webRes.Result {
		return true
	}
	return false
}

// GetWebhookInfo 获取监控信息
func (t *telegramStruct) GetWebhookInfo(Ctx context.Context, Token string) string {
	uri := fmt.Sprintf("%s%s/getWebhookInfo", tgAPI, Token)
	res := g.Client().PostContent(Ctx, uri)
	var webRes webhookInfoRes
	_ = gconv.Struct(res, &webRes)
	if !webRes.Ok {
		_s := gconv.String(webRes.Result)
		g.Log().Stdout(true).Cat("Tg").Print(Ctx, "Webhook:GetWebhookInfo Error"+_s)
		return ""
	}

	return webRes.Result.Url
}

// DelWebhook 删除监控
func (t *telegramStruct) DelWebhook(Ctx context.Context, Token string) bool {
	uri := fmt.Sprintf("%s%s/deleteWebhook", tgAPI, Token)
	var webRes webhookRes
	res := g.Client().PostContent(Ctx, uri)
	_ = gconv.Struct(res, &webRes)
	if !webRes.Ok {
		g.Log().Stdout(true).Cat("Tg").Print(Ctx, "Webhook:DelWebhook Error"+webRes.Description)
		return false
	}
	if webRes.Result {
		return true
	}
	return false
}
