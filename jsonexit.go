package gcmf

import (
	"context"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Exit 退出函数
var Exit = new(exitActive)

type exitActive struct {
}

// JsonResponse 数据返回通用JSON数据结构
type JsonResponse struct {
	Code int         `json:"code"    dc:"Error code"`
	Msg  string      `json:"msg" dc:"Error message"`
	Data interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

// Json 标准返回结果数据结构封装。
func (e *exitActive) Json(r *ghttp.Request, code int, message string, data ...interface{}) {

	var output interface{}
	if len(data) > 0 {
		output = data[0]
	}

	// 自定义空值处理函数

	r.Response.WriteJson(JsonResponse{
		Code: code,
		Msg:  message,
		Data: output,
	})
	r.ExitAll()
}

// API 反回Json格式
func (e *exitActive) API(r *ghttp.Request, data interface{}, Status ...int) {
	_s := 200
	if len(Status) > 0 {
		_s = Status[0]
	}
	r.Response.Status = _s
	r.Response.WriteJson(data)
	r.Exit()
}

// SecondAuth 返回二次验证错误
func (e *exitActive) SecondAuth(Ctx context.Context, authData interface{}) {
	data := JsonResponse{
		Code: 89,
		Msg:  gi18n.T(Ctx, "Need_Second_Auth"),
		Data: authData,
	}
	ghttp.RequestFromCtx(Ctx).Response.WriteJsonExit(data)
}
