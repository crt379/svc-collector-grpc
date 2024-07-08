package appapi

import (
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrappsvc "github.com/crt379/svc-collector-grpc/internal/server/appsvc"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	svrapi "github.com/crt379/svc-collector-grpc/internal/server/svcapi"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.AppapiMeta] = (*AppapiPgDao)(nil)

type AppapiPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *AppapiPgDao) Insert(meta *server.AppapiMeta) (int, error) {
	return 0, fmt.Errorf("method not implemented")
}

func (d *AppapiPgDao) Select(meta *server.AppapiMeta, ops ...server.DaoOption) (objs []server.AppapiMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)
	apid := svrapi.SvcapiPgDao{}
	svcd := svrsvc.ServicePgDao{}
	appsvcd := svrappsvc.AppsvcPgDao{}
	appd := svrapp.ApplicationPgDao{}

	if meta.Appid != 0 {
		k = append(k, d.Field(appsvcd.Table(), "aid"))
		args = append(args, meta.Appid)
	}
	if meta.Appsvcid != 0 {
		k = append(k, d.Field(appsvcd.Table(), "uuid"))
		args = append(args, meta.Appsvcid)
	}
	if meta.Svcid != 0 {
		k = append(k, d.Field(svcd.Table(), "uuid"))
		args = append(args, meta.Svcid)
	}
	if meta.Svcname != "" {
		k = append(k, d.Field(svcd.Table(), "name"))
		args = append(args, meta.Svcname)
	}
	if meta.TenantId != 0 {
		k = append(k, d.Field(svcd.Table(), "tenant_id"))
		args = append(args, meta.TenantId)
	}

	optconditions := make([]string, 0)
	for _, op := range ops {
		optconditions = append(optconditions, op.Conditions()...)
	}
	optconditions = append(optconditions,
		// app_svc_relation.sid = service.uuid
		d.Equal(d.Field(appsvcd.Table(), "sid"), d.Field(svcd.Table(), "uuid")),
		d.Equal(d.Field(appsvcd.Table(), "aid"), d.Field(appd.Table(), "uuid")),
		// service.uuid = service_api.sid
		d.Equal(d.Field(svcd.Table(), "uuid"), d.Field(apid.Table(), "sid")),
	)

	query := d.SelectSQL(
		"",
		d.Comma(svcd.Table(), apid.Table(), appsvcd.Table(), appd.Table()),
		fields(),
		k,
		optconditions...,
	)
	d.Debug(d.Logger, query, args...)

	var rows *sqlx.Rows
	rows, err = d.R.Queryx(query, args...)
	if err != nil {
		return objs, err
	}

	var dbobjs []AppapiDB
	svci := make(map[int]int)
	err = server.RowsToStructs(&dbobjs, rows, func(t *AppapiDB) error {
		i, ok := svci[t.aaservice.Uuid]
		if !ok {
			i = len(objs)
			svci[t.aaservice.Uuid] = i
			objs = append(objs, server.AppapiMeta{
				Appapi: server.AAapi{
					Application: t.ToApplicationMeta(),
					Service:     t.ToServiceMeta(),
					Svcapis:     []server.SvcapiMeta{t.ToSvcapiMeta()},
				},
			})
			return nil
		}

		objs[i].Appapi.Svcapis = append(objs[i].Appapi.Svcapis, t.ToSvcapiMeta())
		return nil
	})

	return objs, err
}

func (d *AppapiPgDao) Count(meta *server.AppapiMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)
	apid := svrapi.SvcapiPgDao{}
	svcd := svrsvc.ServicePgDao{}
	appsvcd := svrappsvc.AppsvcPgDao{}

	if meta.Appid != 0 {
		k = append(k, d.Field(appsvcd.Table(), "aid"))
		args = append(args, meta.Appid)
	}
	if meta.Appsvcid != 0 {
		k = append(k, d.Field(appsvcd.Table(), "uuid"))
		args = append(args, meta.Appsvcid)
	}
	if meta.Svcid != 0 {
		k = append(k, d.Field(svcd.Table(), "uuid"))
		args = append(args, meta.Svcid)
	}
	if meta.Svcname != "" {
		k = append(k, d.Field(svcd.Table(), "name"))
		args = append(args, meta.Svcname)
	}

	query := d.SelectSQL(
		"",
		d.Comma(svcd.Table(), apid.Table(), appsvcd.Table()),
		"count(*)",
		k,
		d.Equal(d.Field(appsvcd.Table(), "sid"), d.Field(svcd.Table(), "uuid")),
		d.Equal(d.Field(svcd.Table(), "uuid"), d.Field(apid.Table(), "sid")),
	)
	d.Debug(d.Logger, query, args...)

	err = d.R.QueryRowx(query, args...).Scan(&count)

	return count, err
}

func (d *AppapiPgDao) Delete(meta *server.AppapiMeta) error {
	return fmt.Errorf("method not implemented")
}

func (d *AppapiPgDao) Update(meta *server.AppapiMeta) (server.AppapiMeta, error) {
	return server.AppapiMeta{}, fmt.Errorf("method not implemented")
}
