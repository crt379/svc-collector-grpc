package service

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.ServiceMeta] = (*ServicePgDao)(nil)

const (
	table = "service"
)

var (
	_fields   = [...]string{"uuid", "name", "describe", "create_time", "update_time", "tenant_id"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type ServicePgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *ServicePgDao) Table() string {
	return table
}

func (d *ServicePgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *ServicePgDao) Insert(meta *server.ServiceMeta) (uuid int, err error) {
	args := []any{meta.Name, meta.Describe, meta.CreateTime, meta.UpdateTime, meta.TenantId}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *ServicePgDao) Select(meta *server.ServiceMeta, ops ...server.DaoOption) (objs []server.ServiceMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Name != "" {
		k = append(k, "name")
		args = append(args, meta.Name)
	}
	if meta.Describe != "" {
		k = append(k, "describe")
		args = append(args, meta.Describe)
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

func (d *ServicePgDao) Count(meta *server.ServiceMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Name != "" {
		k = append(k, "name")
		args = append(args, meta.Name)
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

func (d *ServicePgDao) Delete(meta *server.ServiceMeta) (err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Name != "" {
		k = append(k, "name")
		args = append(args, meta.Name)
	}
	if len(k) == 0 {
		return nil
	}

	query := d.DeleteSQL(d.Table(), k)
	d.Debug(d.Logger, query, args...)

	_, err = d.W.Exec(query, args...)

	return
}

func (d *ServicePgDao) Update(meta *server.ServiceMeta) (obj server.ServiceMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid == 0 {
		return obj, fmt.Errorf("uuid is 0")
	}
	if meta.Name != "" {
		k = append(k, "name")
		args = append(args, meta.Name)
	}
	if meta.Describe != "" {
		k = append(k, "describe")
		args = append(args, meta.Describe)
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
