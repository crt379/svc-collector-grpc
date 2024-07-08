package server

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DaoLog struct{}

func (l *DaoLog) Debug(logger *zap.Logger, query string, args ...any) {
	if logger == nil {
		return
	}
	go func() {
		var buf []byte
		buf, _ = json.Marshal(&args)
		logger.Debug("sql", zap.String("query", query), zap.ByteString("args", buf))
	}()
}

func (l *DaoLog) Info(logger *zap.Logger, query string, args ...any) {
	if logger == nil {
		return
	}
	go func() {
		var buf []byte
		buf, _ = json.Marshal(&args)
		logger.Info("sql", zap.String("query", query), zap.ByteString("args", buf))
	}()
}

type TRow interface {
	any
}

func RowsToStructs[T TRow](objs *[]T, rows *sqlx.Rows, rowfollow ...func(*T) error) (err error) {
	for rows.Next() {
		var row T
		err = rows.StructScan(&row)
		if err != nil {
			return err
		}
		for _, f := range rowfollow {
			err = f(&row)
			if err != nil {
				return err
			}
		}
		(*objs) = append(*objs, row)
	}

	return err
}

type Dao struct{}

func (d *Dao) FieldsP(l int) string {
	if l <= 0 {
		return ""
	}
	switch l {
	case 1:
		return "($1)"
	case 2:
		return "($1, $2)"
	case 3:
		return "($1, $2, $3)"
	case 4:
		return "($1, $2, $3, $4)"
	case 5:
		return "($1, $2, $3, $4, $5)"
	case 6:
		return "($1, $2, $3, $4, $5, $6)"
	case 7:
		return "($1, $2, $3, $4, $5, $6, $7)"
	case 8:
		return "($1, $2, $3, $4, $5, $6, $7, $8)"
	default:
		var fbuilder strings.Builder
		fbuilder.WriteRune('(')
		l := 10
		for i := range l {
			fbuilder.WriteString("$")
			fbuilder.WriteString(strconv.Itoa(i + 1))
			if i != l-1 {
				fbuilder.WriteString(",")
			}
		}
		fbuilder.WriteRune(')')
		return fbuilder.String()
	}
}

func (d *Dao) As(old, new string) string {
	var f strings.Builder
	f.WriteString(old)
	f.WriteString(" as ")
	f.WriteString(new)

	return f.String()
}

func (d *Dao) Field(table, field string) string {
	var f strings.Builder
	f.WriteString(table)
	f.WriteRune('.')
	f.WriteString(field)

	return f.String()
}

func (d *Dao) FieldAs(table string, field string, nfield string) string {
	var f strings.Builder
	f.WriteString(table)
	f.WriteRune('.')
	f.WriteString(field)
	f.WriteString(" as ")
	f.WriteString(nfield)

	return f.String()
}

func (d *Dao) Equal(field1, field2 string) string {
	var f strings.Builder
	f.WriteString(field1)
	f.WriteString(" = ")
	f.WriteString(field2)

	return f.String()
}

func (d *Dao) Comma(field1 string, fields ...string) string {
	if len(fields) == 0 {
		return field1
	}

	var f strings.Builder
	f.WriteString(field1)
	for _, field := range fields {
		f.WriteString(", ")
		f.WriteString(field)
	}

	return f.String()
}

func (d *Dao) FieldIndexP(n string, i int) string {
	var f strings.Builder
	f.WriteString(n)
	f.WriteString("=$")
	f.WriteString(strconv.Itoa(i + 1))

	return f.String()
}

func (d *Dao) WithSQL(wtable string, query string) string {
	ql := []string{
		"WITH", wtable, "AS (",
		query,
		")",
	}

	return strings.Join(ql, " ")
}

func (d *Dao) SelectSQL(with, table string, fields string, cdtfields []string, conditions ...string) string {
	p := 0
	cdtstrs := make([]string, len(cdtfields))
	for i, c := range cdtfields {
		cdtstrs[i] = d.FieldIndexP(c, p)
		p += 1
	}

	if len(cdtstrs) > 0 {
		conditions = append(conditions, cdtstrs...)
	}

	condition := ""
	if len(conditions) > 0 {
		cs := [2]string{"WHERE ", strings.Join(conditions, " AND ")}
		condition = strings.Join(cs[:], "")
	}
	ql := []string{
		with,
		"SELECT", fields,
		"FROM", table,
		condition,
	}

	return strings.Join(ql, " ")
}

func (d *Dao) SelectAddRowNumberSQL(table string, fields string, orderbyField string, cdtfields []string, conditions ...string) string {
	var fbuilder strings.Builder
	fbuilder.WriteString(fields)
	fbuilder.WriteString(", ROW_NUMBER() OVER (ORDER BY ")
	fbuilder.WriteString(orderbyField)
	fbuilder.WriteString(") as row_number")
	return d.SelectSQL("", table, fbuilder.String(), cdtfields, conditions...)
}

func (d *Dao) InsertSQL(table string, fields string, argslen int, returning string) string {
	ret := ""
	if returning != "" {
		ret = "RETURNING " + returning
	}
	ql := []string{
		"INSERT INTO", table,
		"(", fields, ")",
		"VALUES", d.FieldsP(argslen),
		ret,
	}

	return strings.Join(ql, " ")
}

func (d *Dao) UpdateSQL(table string, setfields []string, returning string, cdtfields []string, conditions ...string) string {
	p := 0
	setstrs := make([]string, len(setfields))
	for i, f := range setfields {
		setstrs[i] = d.FieldIndexP(f, p)
		p += 1
	}

	cdtstrs := make([]string, len(cdtfields))
	for i, c := range cdtfields {
		cdtstrs[i] = d.FieldIndexP(c, p)
		p += 1
	}

	if len(cdtstrs) > 0 {
		conditions = append(conditions, cdtstrs...)
	}

	condition := ""
	if len(conditions) > 0 {
		cs := [2]string{"WHERE ", strings.Join(conditions, " AND ")}
		condition = strings.Join(cs[:], "")
	}

	ret := ""
	if returning != "" {
		ret = "RETURNING " + returning
	}

	ql := []string{
		"UPDATE", table,
		"SET", strings.Join(setstrs, ", "),
		condition,
		ret,
	}

	return strings.Join(ql, " ")
}

func (d *Dao) DeleteSQL(table string, cdtfields []string, conditions ...string) string {
	p := 0
	cdtstrs := make([]string, len(cdtfields))
	for i, c := range cdtfields {
		cdtstrs[i] = d.FieldIndexP(c, p)
		p += 1
	}

	if len(cdtstrs) > 0 {
		conditions = append(conditions, cdtstrs...)
	}

	condition := ""
	if len(conditions) > 0 {
		cs := [2]string{"WHERE ", strings.Join(conditions, " AND ")}
		condition = strings.Join(cs[:], "")
	}
	ql := []string{
		"DELETE", "FROM", table,
		condition,
	}

	return strings.Join(ql, " ")
}

type IDao[T any] interface {
	Insert(*T) (int, error)
	Select(*T, ...DaoOption) ([]T, error)
	Count(*T) (int, error)
	Delete(*T) error
	Update(*T) (T, error)
}

type DaoOption interface {
	Conditions() []string
}

type LimitOption struct {
	page  int
	limit int
	field string
}

func NewLimitOption(page int, limit int, field string) *LimitOption {
	return &LimitOption{
		page:  page,
		limit: limit,
		field: field,
	}
}

func (o *LimitOption) Left() string {
	return strconv.Itoa(o.page*o.limit) + " < " + o.field
}

func (o *LimitOption) Right() string {
	return o.field + " <= " + strconv.Itoa((o.page+1)*o.limit)
}

func (o *LimitOption) Conditions() []string {
	if o == nil || o.field == "" || o.limit < 1 {
		return nil
	}
	return []string{o.Left(), o.Right()}
}
