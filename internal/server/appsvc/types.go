package appsvc

import (
	"strings"

	pb "github.com/crt379/svc-collector-grpc-proto/appsvc"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	"github.com/crt379/svc-collector-grpc/internal/types"
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

// type UResp struct {
// 	*pb.UpdateReply
// }

// func (r *UResp) SetCode(code int32) {
// 	r.Code = code
// }

// func (r *UResp) SetMessage(msg string) {
// 	r.Message = msg
// }

// func (r *UResp) GetPBResp() *pb.UpdateReply {
// 	return r.UpdateReply
// }

type appsvc struct {
	Uuid       int        `db:"uuid"`
	AppId      int        `db:"aid"`
	SvcId      int        `db:"sid"`
	CreateTime types.Time `db:"create_time"`
	UpdateTime types.Time `db:"update_time"`
}

type appsvcsvc struct {
	Uuid       int        `db:"s_uuid"`
	Name       string     `db:"s_name"`
	Describe   string     `db:"s_describe"`
	CreateTime types.Time `db:"s_create_time"`
	UpdateTime types.Time `db:"s_update_time"`
	TenantId   int        `db:"s_tenant_id"`
}

var __appsvc_fields string

func fields() string {
	if __appsvc_fields == "" {
		d := AppsvcPgDao{}
		sd := svrsvc.ServicePgDao{}
		fs := []string{
			d.Field(d.Table(), "uuid"),
			d.Field(d.Table(), "aid"),
			d.Field(d.Table(), "sid"),
			d.Field(d.Table(), "create_time"),
			d.Field(d.Table(), "update_time"),
			d.FieldAs(sd.Table(), "uuid", "s_uuid"),
			d.FieldAs(sd.Table(), "name", "s_name"),
			d.FieldAs(sd.Table(), "describe", "s_describe"),
			d.FieldAs(sd.Table(), "create_time", "s_create_time"),
			d.FieldAs(sd.Table(), "update_time", "s_update_time"),
			d.FieldAs(sd.Table(), "tenant_id", "s_tenant_id"),
		}
		__appsvc_fields = strings.Join(fs, ", ")
	}
	return __appsvc_fields
}

type AppsvcSvcDB struct {
	appsvc
	appsvcsvc
}

func (a AppsvcSvcDB) ToAppsvcMeta() server.AppsvcMeta {
	return server.AppsvcMeta{
		Uuid:       a.appsvc.Uuid,
		AppId:      a.appsvc.AppId,
		SvcId:      a.appsvc.SvcId,
		CreateTime: a.appsvc.CreateTime,
		UpdateTime: a.appsvc.UpdateTime,
		Service: server.ServiceMeta{
			Uuid:       a.appsvcsvc.Uuid,
			Name:       a.appsvcsvc.Name,
			Describe:   a.appsvcsvc.Describe,
			CreateTime: a.appsvcsvc.CreateTime,
			UpdateTime: a.appsvcsvc.UpdateTime,
			TenantId:   a.appsvcsvc.TenantId,
		},
	}
}
