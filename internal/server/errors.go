package server

import (
	"github.com/crt379/svc-collector-grpc/internal/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func OkErr() error {
	st := status.New(codes.OK, "ok")

	return st.Err()
}

func InvalidArgumentErr(msg string) error {
	st := status.New(codes.InvalidArgument, "请求参数错误: "+msg)

	return st.Err()
}

func NotFoundErr(msg string) error {
	st := status.New(codes.NotFound, "资源不存在: "+msg)

	return st.Err()
}

func AlreadyExistsErr(msg string) error {
	st := status.New(codes.AlreadyExists, "资源已经存在: "+msg)

	return st.Err()
}

func InternalErr(msg string) error {
	st := status.New(codes.Internal, "服务内部错误: "+msg)

	return st.Err()
}

type PBResp any

type SetResp[T PBResp] interface {
	SetCode(int32)
	SetMessage(string)
	GetMessage() string
	GetPBResp() T
}

func OkResp[T PBResp](resp SetResp[T]) (T, error) {
	resp.SetCode(code.SUCCESS)
	resp.SetMessage("success")

	return resp.GetPBResp(), OkErr()
}

func ParamterResp[T PBResp](resp SetResp[T], msg string) (T, error) {
	resp.SetCode(code.PARAMTER_ERROR)
	resp.SetMessage(msg)

	return resp.GetPBResp(), InvalidArgumentErr(resp.GetMessage())
}

func AlreadyExistsResp[T PBResp](resp SetResp[T], msg string) (T, error) {
	resp.SetCode(code.PARAMTER_ERROR)
	resp.SetMessage(msg)

	return resp.GetPBResp(), AlreadyExistsErr(resp.GetMessage())
}

func SqlErrResp[T PBResp](resp SetResp[T], err error) (T, error) {
	resp.SetCode(code.SQL_EXEC_ERROR)
	resp.SetMessage(err.Error())

	return resp.GetPBResp(), InternalErr(resp.GetMessage())
}

func InternalResp[T PBResp](resp SetResp[T], err error) (T, error) {
	resp.SetCode(code.INTERNAL_ERROR)
	resp.SetMessage(err.Error())

	return resp.GetPBResp(), InternalErr(resp.GetMessage())
}

func NotFoundResp[T PBResp](resp SetResp[T], msg string) (T, error) {
	resp.SetCode(code.PARAMTER_ERROR)
	resp.SetMessage(msg)

	return resp.GetPBResp(), NotFoundErr(resp.GetMessage())
}
