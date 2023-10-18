package valuer

import (
	"database/sql"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/NotFound1911/morm/model"
	"reflect"
)

// reflectValue 基于反射实现
type reflectValue struct {
	val  reflect.Value
	meta *model.Model
}

var _ Creator = NewReflectValue

// NewReflectValue 返回一个封装好的，基于反射实现的 Value
// 输入 val 必须是一个指向结构体实例的指针，而不能是任何其它类型
func NewReflectValue(val interface{}, meta *model.Model) Value {
	return reflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) > len(r.meta.FieldMap) {
		return errs.NewErrTooManyReturnedColumns(cols)
	}
	// colValues 和 colEleValues 实质上最终都指向同一个对象
	colValues := make([]interface{}, len(cols))      // 用于保存每个列
	colEleValues := make([]reflect.Value, len(cols)) // 用于保存每个列的元素值
	for i, col := range cols {
		colFiled, ok := r.meta.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownField(col)
		}
		val := reflect.New(colFiled.Type)
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}
	if err = rows.Scan(colValues...); err != nil {
		return err
	}
	for i, col := range cols {
		colFiled := r.meta.ColumnMap[col]        // 根据列名获取对应的字段结构体
		fd := r.val.FieldByName(colFiled.GoName) // 根据字段名获取字段值
		fd.Set(colEleValues[i])                  // 将查询到的列值设置到字段值上
	}
	return nil
}
