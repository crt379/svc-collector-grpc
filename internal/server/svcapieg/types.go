package svcapieg

import (
	"strings"

	pb "github.com/crt379/svc-collector-grpc-proto/svcapieg"
	"github.com/crt379/svc-collector-grpc/internal/server/jdata"
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

var _svcapi_fields string

func fields() string {
	if _svcapi_fields == "" {
		d := SvcapiegPgDao{}
		jdao := jdata.JdataPgDao{}
		fs := []string{
			d.Field(d.Table(), "uuid"),
			d.Field(jdao.Table(), "data"),
			d.Field(d.Table(), "create_time"),
			d.Field(d.Table(), "update_time"),
			d.Field(d.Table(), "aid"),
			d.Field(d.Table(), "tenant_id"),
			d.Field(d.Table(), "jid"),
		}
		_svcapi_fields = strings.Join(fs, ", ")
	}
	return _svcapi_fields
}
