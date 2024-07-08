package register

import (
	pb "github.com/crt379/svc-collector-grpc-proto/register"
)

type CResp struct {
	*pb.RegisterReply
}

func (r *CResp) SetCode(code int32) {
	r.Code = code
}

func (r *CResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *CResp) GetPBResp() *pb.RegisterReply {
	return r.RegisterReply
}

type GResp struct {
	*pb.GetRegisterReply
}

func (r *GResp) SetCode(code int32) {
	r.Code = code
}

func (r *GResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *GResp) GetPBResp() *pb.GetRegisterReply {
	return r.GetRegisterReply
}

type DResp struct {
	*pb.UnregisterReply
}

func (r *DResp) SetCode(code int32) {
	r.Code = code
}

func (r *DResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *DResp) GetPBResp() *pb.UnregisterReply {
	return r.UnregisterReply
}
