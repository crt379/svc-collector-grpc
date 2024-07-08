package application

import (
	pb "github.com/crt379/svc-collector-grpc-proto/application"
)

type CResp struct {
	*pb.CreateReply
}

func (r *CResp) SetCode(code int32) {
	r.Code = code
}

func (r *CResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *CResp) GetPBResp() *pb.CreateReply {
	return r.CreateReply
}

type GResp struct {
	*pb.GetReply
}

func (r *GResp) SetCode(code int32) {
	r.Code = code
}

func (r *GResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *GResp) GetPBResp() *pb.GetReply {
	return r.GetReply
}

type DResp struct {
	*pb.DeleteReply
}

func (r *DResp) SetCode(code int32) {
	r.Code = code
}

func (r *DResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *DResp) GetPBResp() *pb.DeleteReply {
	return r.DeleteReply
}

type UResp struct {
	*pb.UpdateReply
}

func (r *UResp) SetCode(code int32) {
	r.Code = code
}

func (r *UResp) SetMessage(msg string) {
	r.Message = msg
}

func (r *UResp) GetPBResp() *pb.UpdateReply {
	return r.UpdateReply
}
