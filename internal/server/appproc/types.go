package appproc

import (
	"strings"

	pb "github.com/crt379/svc-collector-grpc-proto/appproc"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrproc "github.com/crt379/svc-collector-grpc/internal/server/processor"
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

type a3app struct {
	Uuid       int        `db:"uuid"`
	Name       string     `db:"name"`
	Describe   string     `db:"describe"`
	CreateTime types.Time `db:"create_time"`
	UpdateTime types.Time `db:"update_time"`
	TenantId   int        `db:"tenant_id"`
}

type a3proc struct {
	Uuid       int        `db:"p_uuid"`
	Addr       string     `db:"p_addr"`
	Weight     int        `db:"p_weight"`
	State      string     `db:"p_state"`
	CreateTime types.Time `db:"p_create_time"`
	UpdateTime types.Time `db:"p_update_time"`
	AppId      int        `db:"p_aid"`
	TanantId   int        `db:"p_tenant_id"`
}

var __fieldstr string

func fields() string {
	if __fieldstr == "" {
		d := server.Dao{}
		appd := svrapp.ApplicationPgDao{}
		procd := svrproc.ProcessorPgDao{}
		fs := []string{
			d.Field(appd.Table(), "uuid"),
			d.Field(appd.Table(), "name"),
			d.Field(appd.Table(), "describe"),
			d.Field(appd.Table(), "create_time"),
			d.Field(appd.Table(), "update_time"),
			d.Field(appd.Table(), "tenant_id"),

			d.FieldAs(procd.Table(), "uuid", "p_uuid"),
			d.FieldAs(procd.Table(), "addr", "p_addr"),
			d.FieldAs(procd.Table(), "weight", "p_weight"),
			d.FieldAs(procd.Table(), "state", "p_state"),
			d.FieldAs(procd.Table(), "create_time", "p_create_time"),
			d.FieldAs(procd.Table(), "update_time", "p_update_time"),
			d.FieldAs(procd.Table(), "aid", "p_aid"),
			d.FieldAs(procd.Table(), "tenant_id", "p_tenant_id"),
		}
		__fieldstr = strings.Join(fs, ", ")
	}

	return __fieldstr
}

type AppprocDB struct {
	a3app
	a3proc
}

func (a AppprocDB) ToApplicationMeta() server.ApplicationMeta {
	return server.ApplicationMeta{
		Uuid:       a.a3app.Uuid,
		Name:       a.a3app.Name,
		Describe:   a.a3app.Describe,
		CreateTime: a.a3app.CreateTime,
		UpdateTime: a.a3app.UpdateTime,
		TenantId:   a.a3app.TenantId,
	}
}

func (a AppprocDB) ToProcessorMeta() server.ProcessorMeta {
	return server.ProcessorMeta{
		Uuid:       a.a3proc.Uuid,
		Addr:       a.a3proc.Addr,
		Weight:     a.a3proc.Weight,
		State:      a.a3proc.State,
		CreateTime: a.a3proc.CreateTime,
		UpdateTime: a.a3proc.UpdateTime,
		AppId:      a.a3proc.AppId,
		TanantId:   a.a3proc.TanantId,
	}
}
