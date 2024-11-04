package gcmf

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type MessageRes struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int `json:"message_id"`
		From      struct {
			Id        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"result"`
}

// SendMessage 获取监控信息
func (t *telegramStruct) SendMessage(Ctx context.Context, Token string, ChatId int64, Text string, ParseMode ...string) bool {
	parseMode := "html"
	if len(ParseMode) > 0 {
		parseMode = ParseMode[0]
	}
	uri := fmt.Sprintf("%s%s/sendMessage", tgAPI, Token)
	_data := g.Map{
		"chat_id":    ChatId,
		"text":       Text,
		"parse_mode": parseMode,
	}
	res := g.Client().PostContent(Ctx, uri, _data)
	var webRes MessageRes
	_ = gconv.Struct(res, &webRes)
	if !webRes.Ok {
		_s := gconv.String(webRes.Result)
		g.Log().Stdout(true).Cat("Tg").Print(Ctx, "Webhook:GetWebhookInfo Error"+_s)
		return false
	}

	return true
}
