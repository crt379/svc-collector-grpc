package appsvc

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.AppsvcMeta] = (*AppsvcPgDao)(nil)

const (
	table = "app_svc_relation"
)

var (
	_fields   = [...]string{"uuid", "aid", "sid", "create_time", "update_time"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type AppsvcPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *AppsvcPgDao) Table() string {
	return table
}

func (d *AppsvcPgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *AppsvcPgDao) Insert(meta *server.AppsvcMeta) (uuid int, err error) {
	args := []any{meta.AppId, meta.SvcId, meta.CreateTime, meta.UpdateTime}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *AppsvcPgDao) Select(meta *server.AppsvcMeta, ops ...server.DaoOption) (objs []server.AppsvcMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.AppId != 0 {
		k = append(k, "aid")
		args = append(args, meta.AppId)
	}
	if meta.SvcId != 0 {
		k = append(k, "sid")
		args = append(args, meta.SvcId)
	}

	optconditions := make([]string, 0)
	for _, op := range ops {
		optconditions = append(optconditions, op.Conditions()...)
	}

	query := d.SelectAddRowNumberSQL(d.Table(), d.fieldsStr(0), "uuid", k)
	query = d.WithSQL("t", query)
	query = d.SelectSQL(query, "t", d.fieldsStr(0), nil, optconditions...)
	d.Debug(d.Logger, query, args...)

	var rows *sqlx.Rows
	rows, err = d.R.Queryx(query, args...)
	if err != nil {
		return objs, err
	}
	err = server.RowsToStructs(&objs, rows)

	return objs, err
}

func (d *AppsvcPgDao) SelectAndService(meta *server.AppsvcMeta, ops ...server.DaoOption) (objs []server.AppsvcMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, d.Field(d.Table(), "uuid"))
		args = append(args, meta.Uuid)
	}
	if meta.AppId != 0 {
		k = append(k, d.Field(d.Table(), "aid"))
		args = append(args, meta.AppId)
	}
	if meta.SvcId != 0 {
		k = append(k, d.Field(d.Table(), "sid"))
		args = append(args, meta.SvcId)
	}

	sd := svrsvc.ServicePgDao{}
	if meta.SvcName != "" {
		k = append(k, d.Field(sd.Table(), "name"))
		args = append(args, meta.SvcName)
	}

	optconditions := make([]string, 0)
	for _, op := range ops {
		optconditions = append(optconditions, op.Conditions()...)
	}
	optconditions = append(optconditions,
		// t1.sid = t2.uuid
		d.Equal(d.Field(d.Table(), "sid"), d.Field(sd.Table(), "uuid")),
	)

	query := d.SelectSQL(
		"",
		d.Comma(d.Table(), sd.Table()),
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

	var dbobjs []AppsvcSvcDB
	err = server.RowsToStructs(&dbobjs, rows, func(t *AppsvcSvcDB) error {
		objs = append(objs, t.ToAppsvcMeta())
		return nil
	})

	return objs, err
}

func (d *AppsvcPgDao) Count(meta *server.AppsvcMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.AppId != 0 {
		k = append(k, "aid")
		args = append(args, meta.AppId)
	}
	if meta.SvcId != 0 {
		k = append(k, "sid")
		args = append(args, meta.SvcId)
	}

	query := d.SelectSQL("", d.Table(), "count(*)", k)
	d.Debug(d.Logger, query, args...)

	err = d.R.QueryRowx(query, args...).Scan(&count)

	return count, err
}

func (d *AppsvcPgDao) Delete(meta *server.AppsvcMeta) (err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if len(k) == 0 {
		return nil
	}

	query := d.DeleteSQL(d.Table(), k)
	d.Debug(d.Logger, query, args...)

	_, err = d.W.Exec(query, args...)

	return
}

func (d *AppsvcPgDao) Update(*server.AppsvcMeta) (server.AppsvcMeta, error) {
	return server.AppsvcMeta{}, fmt.Errorf("method Update not implemented")
}
