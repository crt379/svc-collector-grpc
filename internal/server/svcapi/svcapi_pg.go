package svcapi

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.SvcapiMeta] = (*SvcapiPgDao)(nil)

const (
	table = "service_api"
)

var (
	_fields   = [...]string{"uuid", "path", "method", "describe", "create_time", "update_time", "sid", "tenant_id"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type SvcapiPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *SvcapiPgDao) Table() string {
	return table
}

func (d *SvcapiPgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *SvcapiPgDao) Insert(meta *server.SvcapiMeta) (uuid int, err error) {
	args := []any{meta.Path, meta.Method, meta.Describe, meta.CreateTime, meta.UpdateTime, meta.ServiceId, meta.TenantId}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *SvcapiPgDao) Select(meta *server.SvcapiMeta, ops ...server.DaoOption) (objs []server.SvcapiMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Path != "" {
		k = append(k, "path")
		args = append(args, meta.Path)
	}
	if meta.Method != "" {
		k = append(k, "method")
		args = append(args, meta.Method)
	}
	if meta.Describe != "" {
		k = append(k, "describe")
		args = append(args, meta.Describe)
	}
	if meta.ServiceId != 0 {
		k = append(k, "sid")
		args = append(args, meta.ServiceId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
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

func (d *SvcapiPgDao) Count(meta *server.SvcapiMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Path != "" {
		k = append(k, "path")
		args = append(args, meta.Path)
	}
	if meta.Method != "" {
		k = append(k, "method")
		args = append(args, meta.Method)
	}
	if meta.ServiceId != 0 {
		k = append(k, "sid")
		args = append(args, meta.ServiceId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}

	query := d.SelectSQL("", d.Table(), "count(*)", k)
	d.Debug(d.Logger, query, args...)

	err = d.R.QueryRowx(query, args...).Scan(&count)

	return count, err
}

func (d *SvcapiPgDao) Delete(meta *server.SvcapiMeta) (err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Path != "" {
		k = append(k, "path")
		args = append(args, meta.Path)
	}
	if meta.Method != "" {
		k = append(k, "method")
		args = append(args, meta.Method)
	}
	if meta.ServiceId != 0 {
		k = append(k, "sid")
		args = append(args, meta.ServiceId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
	}
	if len(k) == 0 {
		return nil
	}

	query := d.DeleteSQL(d.Table(), k)
	d.Debug(d.Logger, query, args...)

	_, err = d.W.Exec(query, args...)

	return
}

func (d *SvcapiPgDao) Update(meta *server.SvcapiMeta) (obj server.SvcapiMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid == 0 {
		return obj, fmt.Errorf("uuid is 0")
	}
	if meta.Path != "" {
		k = append(k, "path")
		args = append(args, meta.Path)
	}
	if meta.Method != "" {
		k = append(k, "method")
		args = append(args, meta.Method)
	}
	if meta.Describe != "" {
		k = append(k, "describe")
		args = append(args, meta.Describe)
	}
	if meta.ServiceId != 0 {
		k = append(k, "sid")
		args = append(args, meta.ServiceId)
	}
	if meta.TenantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TenantId)
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
