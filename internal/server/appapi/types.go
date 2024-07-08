package appapi

import (
	"strings"

	pb "github.com/crt379/svc-collector-grpc-proto/appapi"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	svrapi "github.com/crt379/svc-collector-grpc/internal/server/svcapi"
	"github.com/crt379/svc-collector-grpc/internal/types"
)

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

type aaapp struct {
	Uuid       int        `db:"uuid"`
	Name       string     `db:"name"`
	Describe   string     `db:"describe"`
	CreateTime types.Time `db:"create_time"`
	UpdateTime types.Time `db:"update_time"`
	TenantId   int        `db:"tenant_id"`
}

type aaservice struct {
	Uuid       int        `db:"s_uuid"`
	Name       string     `db:"s_name"`
	Describe   string     `db:"s_describe"`
	CreateTime types.Time `db:"s_create_time"`
	UpdateTime types.Time `db:"s_update_time"`
	TenantId   int        `db:"s_tenant_id"`
}

type aasvcapi struct {
	Uuid       int        `db:"a_uuid"`
	Path       string     `db:"a_path"`
	Method     string     `db:"a_method"`
	Describe   string     `db:"a_describe"`
	CreateTime types.Time `db:"a_create_time"`
	UpdateTime types.Time `db:"a_update_time"`
	ServiceId  int        `db:"a_sid"`
	TenantId   int        `db:"a_tenant_id"`
}

var __fieldstr string

func fields() string {
	if __fieldstr == "" {
		d := server.Dao{}
		apid := svrapi.SvcapiPgDao{}
		svcd := svrsvc.ServicePgDao{}
		appd := svrapp.ApplicationPgDao{}
		fs := []string{
			d.Field(appd.Table(), "uuid"),
			d.Field(appd.Table(), "name"),
			d.Field(appd.Table(), "describe"),
			d.Field(appd.Table(), "create_time"),
			d.Field(appd.Table(), "update_time"),
			d.Field(appd.Table(), "tenant_id"),

			d.FieldAs(svcd.Table(), "uuid", "s_uuid"),
			d.FieldAs(svcd.Table(), "name", "s_name"),
			d.FieldAs(svcd.Table(), "describe", "s_describe"),
			d.FieldAs(svcd.Table(), "create_time", "s_create_time"),
			d.FieldAs(svcd.Table(), "update_time", "s_update_time"),
			d.FieldAs(svcd.Table(), "tenant_id", "s_tenant_id"),

			d.FieldAs(apid.Table(), "uuid", "a_uuid"),
			d.FieldAs(apid.Table(), "path", "a_path"),
			d.FieldAs(apid.Table(), "method", "a_method"),
			d.FieldAs(apid.Table(), "describe", "a_describe"),
			d.FieldAs(apid.Table(), "create_time", "a_create_time"),
			d.FieldAs(apid.Table(), "update_time", "a_update_time"),
			d.FieldAs(apid.Table(), "sid", "a_sid"),
			d.FieldAs(apid.Table(), "tenant_id", "a_tenant_id"),
		}
		__fieldstr = strings.Join(fs, ", ")
	}

	return __fieldstr
}

type AppapiDB struct {
	aaapp
	aaservice
	aasvcapi
}

func (a AppapiDB) ToApplicationMeta() server.ApplicationMeta {
	return server.ApplicationMeta{
		Uuid:       a.aaapp.Uuid,
		Name:       a.aaapp.Name,
		Describe:   a.aaapp.Describe,
		CreateTime: a.aaapp.CreateTime,
		UpdateTime: a.aaapp.UpdateTime,
		TenantId:   a.aaapp.TenantId,
	}
}

func (a AppapiDB) ToServiceMeta() server.ServiceMeta {
	return server.ServiceMeta{
		Uuid:       a.aaservice.Uuid,
		Name:       a.aaservice.Name,
		Describe:   a.aaservice.Describe,
		CreateTime: a.aaservice.CreateTime,
		UpdateTime: a.aaservice.UpdateTime,
		TenantId:   a.aaservice.TenantId,
	}
}

func (a AppapiDB) ToSvcapiMeta() server.SvcapiMeta {
	return server.SvcapiMeta{
		Uuid:       a.aasvcapi.Uuid,
		Path:       a.aasvcapi.Path,
		Method:     a.aasvcapi.Method,
		Describe:   a.aasvcapi.Describe,
		CreateTime: a.aasvcapi.CreateTime,
		UpdateTime: a.aasvcapi.UpdateTime,
		ServiceId:  a.aasvcapi.ServiceId,
		TenantId:   a.aasvcapi.TenantId,
	}
}
