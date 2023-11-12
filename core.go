// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package morm

import (
	"context"
	"database/sql"
	"github.com/NotFound1911/morm/internal/valuer"
	"github.com/NotFound1911/morm/model"
)

type core struct {
	r          model.Registry
	dialect    Dialect
	valCreator valuer.Creator
	ms         []Middleware
}

func getMultiHandler[T any](ctx context.Context, c core, sess session, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	tmpls := make([]*T, 0, 0)
	tmpl := new(T)
	meta, err := c.r.Get(tmpl)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	for rows.Next() {
		tmpl := new(T)
		val := c.valCreator(tmpl, meta)
		if err := val.SetColumns(rows); err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		tmpls = append(tmpls, tmpl)
	}
	if len(tmpls) == 0 {
		return &QueryResult{
			Err: sql.ErrNoRows,
		}
	}
	return &QueryResult{
		Result: tmpls,
		Err:    err,
	}
}
func getHandler[T any](ctx context.Context, c core, sess session, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next() {
		return &QueryResult{
			Err: sql.ErrNoRows,
		}
	}
	tmpl := new(T)
	meta, err := c.r.Get(tmpl)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	val := c.valCreator(tmpl, meta)
	err = val.SetColumns(rows)
	return &QueryResult{
		Result: tmpl,
		Err:    err,
	}
}

func get[T any](ctx context.Context, c core, sess session, qc *QueryContext) *QueryResult {
	var handler HanderFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, c, sess, qc)
	}
	ms := c.ms
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	return handler(ctx, qc)
}
func getMulti[T any](ctx context.Context, c core, sess session, qc *QueryContext) *QueryResult {
	var handler HanderFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getMultiHandler[T](ctx, c, sess, qc)
	}
	ms := c.ms
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	return handler(ctx, qc)
}
func exec(ctx context.Context, sess session, c core, qc *QueryContext) Result {
	var handler HanderFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		res, err := sess.execContext(ctx, q.SQL, q.Args...)
		return &QueryResult{Err: err, Result: res}
	}
	ms := c.ms
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	qr := handler(ctx, qc)
	var res sql.Result
	if qr.Result != nil {
		res = qr.Result.(sql.Result)
	}
	return Result{
		err: qr.Err,
		res: res,
	}
}
