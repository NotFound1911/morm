package morm

import (
	"context"
	"github.com/NotFound1911/morm/errors"
	"reflect"
)

type Updater[T any] struct {
	builder
	val     *T
	assigns []Assignable
	sess    session
}

func NewUpdater[T any](sess session) *Updater[T] {
	c := sess.getCore()
	return &Updater[T]{
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
			quoter:  c.dialect.quoter(),
		},
	}
}

func (u *Updater[T]) Update(t *T) *Updater[T] {
	u.val = t
	return u
}

func (u *Updater[T]) Set(assigns ...Assignable) *Updater[T] {
	u.assigns = assigns
	return u
}

func (u *Updater[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	u.model, err = u.r.Get(&t)
	if err != nil {
		return nil, err
	}
	if len(u.assigns) == 0 {
		return nil, errs.NewErrNoUpdatedColumns()
	}
	u.sqlBuilder.WriteString("UPDATE ")
	if err = u.buildTable(u.table); err != nil {
		return nil, err
	}
	u.sqlBuilder.WriteString(" SET ")
	val := u.valCreator(u.val, u.model)
	for i := 0; i < len(u.assigns); i++ {
		if i > 0 {
			u.sqlBuilder.WriteByte(',')
		}
		switch assign := u.assigns[i].(type) {
		case Column:
			if err := u.buildColumn(assign.table, assign.name); err != nil {
				return nil, err
			}
			u.sqlBuilder.WriteString("=?")
			arg, err := val.Field(assign.name)
			if err != nil {
				return nil, err
			}
			u.addArgs(arg)
		case Assignment:
			if err := u.buildAssignment(assign); err != nil {
				return nil, err
			}
		default:
			return nil, errs.NewErrUnsupportedAssignableType(assign)

		}
	}
	if len(u.where) > 0 {
		u.sqlBuilder.WriteString(" WHERE ")
		if err := u.buildPredicates(u.where); err != nil {
			return nil, err
		}
	}
	u.sqlBuilder.WriteByte(';')
	return &Query{
		SQL:  u.sqlBuilder.String(),
		Args: u.args,
	}, nil
}
func (u *Updater[T]) buildTable(table TableReference) error {
	switch tab := table.(type) {
	case nil:
		u.quote(u.model.TableName)
	case Table:
		model, err := u.r.Get(tab.entity)
		if err != nil {
			return err
		}
		u.quote(model.TableName)
		if tab.alias != "" {
			u.sqlBuilder.WriteString(" AS ")
			u.quote(tab.alias)
		}
	default:
		return errs.NewErrUnsupportedExpressionType(tab)
	}
	return nil
}
func (u *Updater[T]) buildAssignment(assign Assignment) error {
	if err := u.buildColumn(nil, assign.Column.name); err != nil {
		return err
	}
	u.sqlBuilder.WriteByte('=')
	return u.buildExpression(assign.val)
}
func (u *Updater[T]) Where(ps ...Predicate) *Updater[T] {
	u.where = ps
	return u
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	return exec(ctx, u.sess, u.core, &QueryContext{Builder: u, Type: "UPDATE"})
}

func AssignNotNilColumns(entity interface{}) []Assignable {
	return AssignColumns(entity, func(typ reflect.StructField, val reflect.Value) bool {
		switch val.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			return !val.IsNil()
		}
		return true
	})
}

func AssignNotZeroColumns(entity interface{}) []Assignable {
	return AssignColumns(entity, func(typ reflect.StructField, val reflect.Value) bool {
		return !val.IsZero()
	})
}
func AssignColumns(entity interface{}, filter func(typ reflect.StructField, val reflect.Value) bool) []Assignable {
	val := reflect.ValueOf(entity).Elem()
	typ := reflect.TypeOf(entity).Elem()
	numField := val.NumField()
	res := make([]Assignable, 0, numField)
	for i := 0; i < numField; i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		if filter(fieldTyp, fieldVal) {
			res = append(res, Assign(fieldTyp.Name, fieldVal.Interface()))
		}
	}
	return res
}
