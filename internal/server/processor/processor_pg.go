package processor

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.ProcessorMeta] = (*ProcessorPgDao)(nil)

const (
	table = "processor"
)

var (
	_fields   = [...]string{"uuid", "addr", "weight", "state", "create_time", "update_time", "aid", "tenant_id"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type ProcessorPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *ProcessorPgDao) Table() string {
	return table
}

func (d *ProcessorPgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *ProcessorPgDao) Insert(meta *server.ProcessorMeta) (uuid int, err error) {
	args := []any{meta.Addr, meta.Weight, meta.State, meta.CreateTime, meta.UpdateTime, meta.AppId, meta.TanantId}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *ProcessorPgDao) Select(meta *server.ProcessorMeta, ops ...server.DaoOption) (objs []server.ProcessorMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Addr != "" {
		k = append(k, "addr")
		args = append(args, meta.Addr)
	}
	if meta.Weight != 0 {
		k = append(k, "weight")
		args = append(args, meta.Weight)
	}
	if meta.State != "" {
		k = append(k, "state")
		args = append(args, meta.State)
	}
	if meta.AppId != 0 {
		k = append(k, "aid")
		args = append(args, meta.AppId)
	}
	if meta.TanantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TanantId)
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

func (d *ProcessorPgDao) Count(meta *server.ProcessorMeta) (count int, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.Addr != "" {
		k = append(k, "addr")
		args = append(args, meta.Addr)
	}
	if meta.Weight != 0 {
		k = append(k, "weight")
		args = append(args, meta.Weight)
	}
	if meta.State != "" {
		k = append(k, "state")
		args = append(args, meta.State)
	}
	if meta.AppId != 0 {
		k = append(k, "aid")
		args = append(args, meta.AppId)
	}
	if meta.TanantId != 0 {
		k = append(k, "tenant_id")
		args = append(args, meta.TanantId)
	}

	query := d.SelectSQL("", d.Table(), "count(*)", k)
	d.Debug(d.Logger, query, args...)

	err = d.R.QueryRowx(query, args...).Scan(&count)

	return count, err
}

func (d *ProcessorPgDao) Delete(meta *server.ProcessorMeta) (err error) {
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

func (d *ProcessorPgDao) Update(meta *server.ProcessorMeta) (obj server.ProcessorMeta, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid == 0 {
		return obj, fmt.Errorf("uuid is 0")
	}

	if meta.Addr != "" {
		k = append(k, "addr")
		args = append(args, meta.Addr)
	}
	if meta.Weight != 0 {
		k = append(k, "weight")
		args = append(args, meta.Weight)
	}
	if meta.State != "" {
		k = append(k, "state")
		args = append(args, meta.State)
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
