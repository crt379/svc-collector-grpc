package jdata

import (
	"fmt"
	"strings"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ server.IDao[server.Jdata] = (*JdataPgDao)(nil)

const (
	table = "jdata"
)

var (
	_fields   = [...]string{"uuid", "data", "create_time", "update_time", "hash_type", "hash_value"}
	_fields_0 = strings.Join(_fields[:], ",")
	_fields_1 = strings.Join(_fields[1:], ",")
)

type JdataPgDao struct {
	W      *sqlx.DB
	R      *sqlx.DB
	Logger *zap.Logger
	server.Dao
	server.DaoLog
}

func (d *JdataPgDao) Table() string {
	return table
}

func (d *JdataPgDao) fieldsStr(s int) string {
	switch s {
	case 0:
		return _fields_0
	case 1:
		return _fields_1
	}
	return strings.Join(_fields[s:], ",")
}

func (d *JdataPgDao) Insert(meta *server.Jdata) (uuid int, err error) {
	args := []any{meta.Data, meta.CreateTime, meta.UpdateTime, meta.HashType, meta.HashValue}

	query := d.InsertSQL(d.Table(), d.fieldsStr(1), len(args), "uuid")
	d.Debug(d.Logger, query, args...)

	err = d.W.QueryRowx(query, args...).Scan(&uuid)

	return uuid, err
}

func (d *JdataPgDao) Select(meta *server.Jdata, ops ...server.DaoOption) (objs []server.Jdata, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.HashType != "" {
		k = append(k, "hash_type")
		args = append(args, meta.HashType)
	}
	if meta.HashValue != "" {
		k = append(k, "hash_value")
		args = append(args, meta.HashValue)
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

func (d *JdataPgDao) Count(*server.Jdata) (int, error) {
	return 0, fmt.Errorf("method Count not implemented")
}

func (d *JdataPgDao) Delete(meta *server.Jdata) (err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid != 0 {
		k = append(k, "uuid")
		args = append(args, meta.Uuid)
	}
	if meta.HashType != "" {
		k = append(k, "hash_type")
		args = append(args, meta.HashType)
	}
	if meta.HashValue != "" {
		k = append(k, "hash_value")
		args = append(args, meta.HashValue)
	}
	if len(k) == 0 {
		return nil
	}

	query := d.DeleteSQL(d.Table(), k)
	d.Debug(d.Logger, query, args...)

	_, err = d.W.Exec(query, args...)

	return
}

func (d *JdataPgDao) Update(meta *server.Jdata) (obj server.Jdata, err error) {
	k := make([]string, 0)
	args := make([]any, 0)

	if meta.Uuid == 0 {
		return obj, fmt.Errorf("uuid is 0")
	}

	k = append(k, "data")
	args = append(args, meta.Data)

	if meta.HashType != "" {
		k = append(k, "hash_type")
		args = append(args, meta.HashType)
	}
	if meta.HashValue != "" {
		k = append(k, "hash_value")
		args = append(args, meta.HashValue)
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
