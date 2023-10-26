package morm

type Column struct {
	name  string
	alias string
}

func (c Column) selectable() {
}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

var _ Expression = &Column{}

func (c Column) expr() {

}

type value struct {
	val any
}

var _ Expression = &value{}

func (v value) expr() {

}

func valueOf(val any) value {
	return value{
		val: val,
	}
}

func C(name string) Column {
	return Column{
		name: name,
	}
}

//  对应列的方法，= < >

// EQ C("id").EQ(12)
func (c Column) EQ(arg any) Predicate { // = 等于
	return Predicate{
		left:  c,
		opt:   optEQ,
		right: exprOf(arg),
	}
}
func (c Column) LT(arg any) Predicate { // < 小于
	return Predicate{
		left:  c,
		opt:   optLT,
		right: exprOf(arg),
	}
}
func (c Column) GT(arg any) Predicate { // > 大于
	return Predicate{
		left:  c,
		opt:   optGT,
		right: exprOf(arg),
	}
}
