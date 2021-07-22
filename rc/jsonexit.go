package rc

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

// JsonResponse 数据返回通用JSON数据结构
type JsonResponse struct {
	Code int         `json:"code"` // 错误码((0:失败, >0:错误码))
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 返回数据(业务接口定义具体数据结构)
}

// LayuiResponse layui 返回数据
type LayuiResponse struct {
	Code  int         `json:"code"`  // 错误码((0:失败, >0:错误码))
	Msg   string      `json:"msg"`   // 提示信息
	Data  interface{} `json:"data"`  // 返回数据(业务接口定义具体数据结构)
	Count int         `json:"count"` // 提示信息
}

// Params is type for template params.

// Json 标准返回结果数据结构封装。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}

	_ = r.Response.WriteJson(JsonResponse{
		Code: code,
		Msg:  message,
		Data: responseData,
	})
}

// JsonExit 返回JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, err int, msg string, data ...interface{}) {
	Json(r, err, msg, data...)
	r.Exit()
}

// LayuiData 数据
func LayuiData(r *ghttp.Request, code int, msg string, count int, data interface{}) {

	_ = r.Response.WriteJson(LayuiResponse{
		Code:  code,
		Msg:   msg,
		Count: count,
		Data:  data,
	})
	r.Exit()
}

// DataError
func DataError(r *ghttp.Request, msg string) {
	//	html := "<script type=\"text/javascript\">layui.admin.rc_msg(\"" + msg + "\",1)</script>"
	//if r.IsAjaxRequest() {
	r.Response.WriteStatus(http.StatusNotImplemented, msg)
	//r.Response.Write(html)
	r.Exit()
	//}
	// 跳转页面
}

// APIExit 反回Json格式
func APIExit(r *ghttp.Request, data interface{}, Status ...int) {
	_s := 200
	if len(Status) > 0 {
		_s = Status[0]
	}
	r.Response.Status = _s
	_ = r.Response.WriteJson(data)
	r.Exit()
}
