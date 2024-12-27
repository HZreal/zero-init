package xerr

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
常用通用固定错误
*/

type CodeError struct {
	errCode         uint32
	errMsg          string
	errDetail       string
	errMsgParams    []string
	errDetailParams []string
}

// 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

// 返回给前端显示端详细错误信息
func (e *CodeError) GetErrDetail() string {
	return e.errDetail
}

// 返回给前端Msg错误信息中的参数
func (e *CodeError) GetErrMsgParams() []string {
	return e.errMsgParams
}

// 返回给前端Detail错误信息中的参数
func (e *CodeError) GetErrDetailParams() []string {
	return e.errDetailParams
}

func (e *CodeError) SetErrMsgParams(params []string) {
	e.errMsgParams = params
}

func (e *CodeError) SetErrDetailParams(params []string) {
	e.errDetailParams = params
}

func (e *CodeError) SetErrDetail(detail string) {
	e.errDetail = detail
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

// func (e *CodeError) Is(target error) bool {
// 	// 判断目标是否是 CodeError 类型
// 	t, ok := target.(*CodeError)
// 	// var t *CodeError
// 	// ok := errors.As(target, &t)
// 	if !ok {
// 		return false
// 	}
// 	// 比较错误码
// 	return e.errCode == t.errCode
// }

func NewErrCodeMsgParams(errCode uint32, errMsg string, params []string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg, errMsgParams: params}
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ServerBusy, errMsg: errMsg}
}

// 把异常转化为rpc的error
func NewRpcErr(e *CodeError) error {
	return status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
}

// 构造一个rpc error。
func NewRpcErrCode(errCode uint32) error {
	err := &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
	return NewRpcErr(err)
}

// 构造一个带参数的rpc错误
func NewRpcErrCodeParams(errCode uint32, params ...any) error {
	return status.Errorf(codes.Code(errCode), MapErrMsg(errCode), params...)
}
