package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/NotFound1911/morm/model"
	"reflect"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []Column
}

type Upsert struct {
	conflictColumns []Column
	assigns         []Assignable
}

func (u *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	u.conflictColumns = make([]Column, len(cols))
	for i, col := range cols {
		u.conflictColumns[i] = Column{
			name: col,
		}
	}
	return u
}
func (u *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	u.i.onDuplicate = &Upsert{
		conflictColumns: u.conflictColumns,
		assigns:         assigns,
	}
	return u.i
}

type Inserter[T any] struct {
	builder
	values      []*T     // 插入值
	columns     []string // 指定列
	onDuplicate *Upsert
}

// OnDuplicateKey  返回OnDuplicateKey构造部分
// 整体为 Inserter构造 --> OnDuplicateKey构造冲突部分 --> Inserter构造剩余部分
func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}
func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: builder{
			db:      db,
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
	}
}

// Values 要插入的值
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) Cloumns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}
func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.NewErrInsertZeroRow()
	}
	var (
		t   T
		err error
	)
	i.model, err = i.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	i.sqlBuilder.WriteString("INSERT INTO `")
	i.sqlBuilder.WriteString(i.model.TableName)
	i.sqlBuilder.WriteString("`(")

	fields := i.model.Fields
	if len(i.columns) != 0 { // 指定列
		fields = make([]*model.Field, 0, len(i.columns))
		for _, col := range i.columns { // 使用sql的顺序
			field, ok := i.model.FieldMap[col]
			if !ok {
				return nil, errs.NewErrUnknownField(col)
			}
			fields = append(fields, field)
		}
	}
	// (len(i.values) + 1) 中 +1 是考虑到 UPSERT 语句会传递额外的参数
	i.args = make([]any, 0, len(fields)*(len(i.values)+1))
	for idx, fd := range fields {
		if idx > 0 {
			i.sqlBuilder.WriteByte(',')
		}
		i.sqlBuilder.WriteByte('`')
		i.sqlBuilder.WriteString(fd.ColName)
		i.sqlBuilder.WriteByte('`')
	}
	i.sqlBuilder.WriteString(") VALUES")
	for vIdx, val := range i.values { // 第一层便利值
		if vIdx > 0 {
			i.sqlBuilder.WriteByte(',')
		}
		refVal := reflect.ValueOf(val).Elem()
		i.sqlBuilder.WriteByte('(')
		for fIdx, field := range fields { // 第二层便利字段
			if fIdx > 0 {
				i.sqlBuilder.WriteByte(',')
			}
			i.sqlBuilder.WriteByte('?')
			fdVal := refVal.Field(field.Index)
			i.addArgs(fdVal.Interface())
		}
		i.sqlBuilder.WriteByte(')')
	}
	// 构造冲突部分
	if i.onDuplicate != nil {
		if err := i.dialect.buildUpsert(&i.builder, i.onDuplicate); err != nil {
			return nil, err
		}
	}
	i.sqlBuilder.WriteByte(';')
	return &Query{
		SQL:  i.sqlBuilder.String(),
		Args: i.args,
	}, nil
}
func (i *Inserter[T]) buildAssignment(a Assignable) error {
	switch assign := a.(type) {
	case Column:
		i.sqlBuilder.WriteByte('`')
		fd, ok := i.model.FieldMap[assign.name]
		if !ok {
			return errs.NewErrUnknownField(assign.name)
		}
		i.sqlBuilder.WriteString(fd.ColName)
		i.sqlBuilder.WriteString("`=VALUES(`")
		i.sqlBuilder.WriteString(fd.ColName)
		i.sqlBuilder.WriteString("`)")
	case Assignment:
		i.sqlBuilder.WriteByte('`')
		fd, ok := i.model.FieldMap[assign.name]
		if !ok {
			return errs.NewErrUnknownField(assign.name)
		}
		i.sqlBuilder.WriteString(fd.ColName)
		i.sqlBuilder.WriteByte('`')
		i.sqlBuilder.WriteString("=?")
		i.addArgs(assign.val)
	default:
		return errs.NewErrUnsupportedAssignableType(a)
	}
	return nil
}
