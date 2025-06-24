package response

import (
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"reflect"
)

type JsonResponse struct {
	Code    int32  `json:"code"`
	Error   any    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func http(r flamego.Render, code int32, msg string, data any, err any) {
	if code == 0 {
		r.JSON(200, &JsonResponse{
			Code:    code,
			Message: msg,
			Data:    data,
		})
		return
	}
	r.JSON(int(code/1000), &JsonResponse{
		Code:    code,
		Message: msg,
		Data:    data,
		Error:   err,
	})
}

// HTTPSuccess 成功返回
func HTTPSuccess(r flamego.Render, data any) {
	http(r, 0, "success", data, nil)
}

func HTTPFail(r flamego.Render, code int, msg string, err ...any) {
	newErrs := make([]string, 0, len(err))
	for _, e := range err {
		if reflect.TypeOf(e).Kind() == reflect.Slice {
			if errs, ok := e.([]error); ok {
				for _, e := range errs {
					newErrs = append(newErrs, e.Error())
				}
				continue
			}
			if errs, ok := e.(binding.Errors); ok {
				for _, e := range errs {
					newErrs = append(newErrs, e.Err.Error())
				}
				continue
			}
		} else {
			if v, ok := e.(error); ok {
				newErrs = append(newErrs, v.Error())
				continue
			}
			if v, ok := e.(string); ok {
				newErrs = append(newErrs, v)
				continue
			}

		}
	}
	http(r, int32(code), msg, nil, newErrs)
}

func HTTPFailWithData(r flamego.Render, code int, msg string, data any, err ...any) {
	for i, e := range err {
		if v, ok := e.(error); ok {
			err[i] = v.Error()
		}
	}
	http(r, int32(code), msg, data, err)
}

func UnAuthorization(r flamego.Render) {
	HTTPFail(r, 401000, "登录过期失效，请重新登陆")
}

func Forbidden(r flamego.Render) {
	HTTPFail(r, 403000, "未授权")
}

func InValidParam(r flamego.Render, err ...any) {
	HTTPFail(r, 400000, "请求校验失败", err...)
}

func ServiceErr(r flamego.Render, err ...any) {
	HTTPFail(r, 400000, "内部异常", err...)
}
