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
	"github.com/NotFound1911/morm/model"
)

// Value 是结构体实例在内存中的具体数据表示
type Value interface {
	// Field 返回字段对应的值
	Field(name string) (any, error)
	// SetColumns 设置列值
	SetColumns(rows *sql.Rows) error
}
type Creator func(val interface{}, meta *model.Model) Value
