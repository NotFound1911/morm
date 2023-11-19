package morm

// Assignable 标记接口
// 可用于UPDATE和UPSERT 中赋值语句
type Assignable interface {
	assign()
}
type Assignment struct {
	Column
	val Expression
}

func (a Assignment) assign() {

}

func Assign(column string, val any) Assignment {
	v, ok := val.(Expression)
	if !ok {
		v = value{val: val}
	}
	return Assignment{
		Column: Column{
			name: column,
		},
		val: v,
	}
}
