package svcapieg

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/server/jdata"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.SvcapiegMeta] = (*SvcapiegPgDao)(nil)

const (
	table = "svc_api_example"
)

var (
	_fields   = [...]string{"uuid", "create_time", "update_time", "aid", "tenant_id", "jid"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type SvcapiegPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *SvcapiegPgDao) Table() string {
	return table
}

func (d *SvcapiegPgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *SvcapiegPgDao) Insert(meta *server.SvcapiegMeta) (uuid int, err error) {
	args := []any{meta.CreateTime, meta.UpdateTime, meta.SvcapiId, meta.TenantId, meta.JdataId}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *SvcapiegPgDao) Select(meta *server.SvcapiegMeta, ops ...server.DaoOption) (objs []server.SvcapiegMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.SvcapiId != 0 {
		k = append(k, "aid")
		args = append(args, meta.SvcapiId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}
	if meta.JdataId != 0 {
		k = append(k, "jid")
		args = append(args, meta.JdataId)
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

func (d *SvcapiegPgDao) SelectAndJdata(meta *server.SvcapiegMeta, ops ...server.DaoOption) (objs []server.SvcapiegMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, d.Field(d.Table(), "uuid"))
		args = append(args, meta.Uuid)
	}
	if meta.SvcapiId != 0 {
		k = append(k, d.Field(d.Table(), "aid"))
		args = append(args, meta.SvcapiId)
	}
	if meta.TenantId != 0 {
		k = append(k, d.Field(d.Table(), "tenant_id"))
		args = append(args, meta.TenantId)
	}
	if meta.JdataId != 0 {
		k = append(k, d.Field(d.Table(), "jid"))
		args = append(args, meta.JdataId)
	}

	jdao := jdata.JdataPgDao{}

	query := d.SelectSQL(
		"",
		d.Comma(d.Table(), jdao.Table()),
		fields(),
		k,
		d.Equal(d.Field(d.Table(), "jid"), d.Field(jdao.Table(), "uuid")),
	)
	d.Debug(d.Logger, query, args...)

	var rows *sqlx.Rows
	rows, err = d.R.Queryx(query, args...)
	if err != nil {
		return objs, err
	}
	err = server.RowsToStructs(&objs, rows)

	return objs, err
}

func (d *SvcapiegPgDao) Count(meta *server.SvcapiegMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.SvcapiId != 0 {
		k = append(k, "aid")
		args = append(args, meta.SvcapiId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}
	if meta.JdataId != 0 {
		k = append(k, "jid")
		args = append(args, meta.JdataId)
	}

	query := d.SelectSQL("", d.Table(), "count(*)", k)
	d.Debug(d.Logger, query, args...)

	err = d.R.QueryRowx(query, args...).Scan(&count)

	return count, err
}

func (d *SvcapiegPgDao) Delete(meta *server.SvcapiegMeta) (err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.SvcapiId != 0 {
		k = append(k, "aid")
		args = append(args, meta.SvcapiId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}
	if meta.JdataId != 0 {
		k = append(k, "jid")
		args = append(args, meta.JdataId)
	}
	if len(k) == 0 {
		return nil
	}

	query := d.DeleteSQL(d.Table(), k)
	d.Debug(d.Logger, query, args...)

	_, err = d.W.Exec(query, args...)

	return
}

func (d *SvcapiegPgDao) Update(meta *server.SvcapiegMeta) (obj server.SvcapiegMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid == 0 {
		return obj, fmt.Errorf("uuid is 0")
	}
	if meta.SvcapiId != 0 {
		k = append(k, "aid")
		args = append(args, meta.SvcapiId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}
	if meta.JdataId != 0 {
		k = append(k, "jid")
		args = append(args, meta.JdataId)
	}

	var zerotime types.Time
	if meta.UpdateTime != zerotime {
		k = append(k, "update_time")
		args = append(args, meta.UpdateTime)
	}

	args = append(args, meta.Uuid)
	query := d.UpdateSQL(d.Table(), k, d.fieldsStr(0), []string{"uuid"})
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).StructScan(&obj)

	return obj, err
}
