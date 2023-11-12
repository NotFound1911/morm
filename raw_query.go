package morm

import "context"

var _ Querier[any] = &RawQuerier[any]{}

type RawQuerier[T any] struct {
	core
	sess session
	sql  string
	args []any
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func (r *RawQuerier[T]) Exec(ctx context.Context) Result {
	return exec(ctx, r.sess, r.core, &QueryContext{
		Builder: r,
		Type:    "RAW",
	})
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	res := get[T](ctx, r.core, r.sess, &QueryContext{
		Builder: r,
		Type:    "RAW",
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	res := getMulti[T](ctx, r.core, r.sess, &QueryContext{
		Builder: r,
		Type:    "RAW",
	})
	if res.Result != nil {
		return res.Result.([]*T), res.Err
	}
	return nil, res.Err
}

func RawQuery[T any](sess session, sql string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  sql,
		args: args,
		core: sess.getCore(),
		sess: sess,
	}
}
