package response

import (
	"overall/common/xerr"
	"reflect"
)

type ResultType struct {
	Code   uint32      `json:"code"`
	Msg    string      `json:"msg"`
	Detail string      `json:"detail"`
	Data   interface{} `json:"data"`
}

type NullJson struct{}

func Success(data interface{}) *ResultType {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind().String() == "struct" {
		if v.IsValid() {
			v11 := v.FieldByName("Data")
			if v11.IsValid() {
				data = v11.Interface()
			}
		}
	}

	return &ResultType{
		Code: xerr.OK,
		Msg:  xerr.MapErrMsg(xerr.OK),
		// Msg:  i18n.TS(xerr.MapErrMsg(xerr.OK), ""),
		Data: data,
	}
}

func Error(errCode uint32, errMsg string, errDetail string, data interface{}) (result *ResultType) {
	defer func() {
		if err := recover(); err != nil {
			result = &ResultType{
				Code:   errCode,
				Msg:    errMsg,
				Detail: errDetail,
				Data:   nil,
			}
		}
	}()

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind().String() == "struct" {
		if v.IsValid() {
			v11 := v.FieldByName("Data")
			if v11.IsValid() {
				data = v11.Interface()
			}
		}
	}

	result = &ResultType{
		Code: errCode,
		Msg:  errMsg,
		Data: data,
	}
	return result
}
