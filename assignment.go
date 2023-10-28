package morm

// Assignable 标记接口
// 可用于UPDATE和UPSERT 中赋值语句
type Assignable interface {
	assign()
}
type Assignment struct {
	column string
	val    any
}

func (a Assignment) assign() {

}

func Assign(column string, val any) Assignment {
	return Assignment{
		column: column,
		val:    val,
	}
}
