package appproc

import (
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrproc "github.com/crt379/svc-collector-grpc/internal/server/processor"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.AppprocMeta] = (*AppprocPgDao)(nil)

type AppprocPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *AppprocPgDao) Insert(meta *server.AppprocMeta) (int, error) {
	return 0, fmt.Errorf("method not implemented")
}

func (d *AppprocPgDao) Select(meta *server.AppprocMeta, ops ...server.DaoOption) (objs []server.AppprocMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	appd := svrapp.ApplicationPgDao{}
	procd := svrproc.ProcessorPgDao{}

	if meta.Appid != 0 {
		k = append(k, d.Field(appd.Table(), "uuid"))
		args = append(args, meta.Appid)
	}
	if meta.Appname != "" {
		k = append(k, d.Field(appd.Table(), "name"))
		args = append(args, meta.Appname)
	}
	if meta.Weight != nil {
		k = append(k, d.Field(procd.Table(), "weight"))
		args = append(args, *meta.Weight)
	}
	if meta.State != "" {
		k = append(k, d.Field(procd.Table(), "state"))
		args = append(args, meta.State)
	}
	if meta.TenantId != 0 {
		k = append(k, d.Field(appd.Table(), "tenant_id"))
		args = append(args, meta.TenantId)
	}

	optconditions := make([]string, 0)
	for _, op := range ops {
		optconditions = append(optconditions, op.Conditions()...)
	}
	optconditions = append(optconditions,
		d.Equal(d.Field(appd.Table(), "uuid"), d.Field(procd.Table(), "aid")),
	)

	query := d.SelectSQL(
		"",
		d.Comma(procd.Table(), appd.Table()),
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

	var dbobjs []AppprocDB
	appi := make(map[int]int)
	err = server.RowsToStructs(&dbobjs, rows, func(t *AppprocDB) error {
		i, ok := appi[t.a3app.Uuid]
		if !ok {
			i = len(objs)
			appi[t.a3app.Uuid] = i
			objs = append(objs, server.AppprocMeta{
				A3p: server.A3p{
					Application: t.ToApplicationMeta(),
					Processors:  []server.ProcessorMeta{t.ToProcessorMeta()},
				},
			})
			return nil
		}

		objs[i].A3p.Processors = append(objs[i].A3p.Processors, t.ToProcessorMeta())
		return nil
	})

	return objs, err
}

func (d *AppprocPgDao) Count(meta *server.AppprocMeta) (count int, err error) {
	return 0, fmt.Errorf("method not implemented")
}

func (d *AppprocPgDao) Delete(meta *server.AppprocMeta) error {
	return fmt.Errorf("method not implemented")
}

func (d *AppprocPgDao) Update(meta *server.AppprocMeta) (server.AppprocMeta, error) {
	return server.AppprocMeta{}, fmt.Errorf("method not implemented")
}
