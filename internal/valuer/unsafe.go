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

package valuer

import (
	"database/sql"
	"github.com/NotFound1911/morm/errors"
	"github.com/NotFound1911/morm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	addr unsafe.Pointer
	meta *model.Model
}

func (u unsafeValue) Field(name string) (any, error) {
	fd, ok := u.meta.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	ptr := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
	val := reflect.NewAt(fd.Type, ptr).Elem()
	return val.Interface(), nil
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(val interface{}, meta *model.Model) Value {
	return unsafeValue{
		addr: unsafe.Pointer(reflect.ValueOf(val).Pointer()),
		meta: meta,
	}
}

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) > len(u.meta.ColumnMap) {
		return errs.NewErrTooManyReturnedColumns(cols)
	}
	colValues := make([]interface{}, len(cols))
	for i, col := range cols {
		colField, ok := u.meta.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownField(col)
		}
		ptr := unsafe.Pointer(uintptr(u.addr) + colField.Offset)
		val := reflect.NewAt(colField.Type, ptr)
		colValues[i] = val.Interface()
	}
	return rows.Scan(colValues...)
}
