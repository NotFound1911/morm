package morm

type Column struct {
	table TableReference
	name  string
	alias string
}

func (c Column) fieldName() string {
	return c.name
}

func (c Column) target() TableReference {
	return c.table
}

func (c Column) selectedAlias() string {
	return c.alias
}

func (c Column) tableAlias() string {
	return c.alias
}

func (c Column) assign() {}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

var _ Expression = &Column{}

func (c Column) expr() {}

type value struct {
	val any
}

var _ Expression = &value{}

func (v value) expr() {}

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

func (c Column) Add(delta int) MathExpr {
	return MathExpr{
		left: c,
		opt:  optADD,
		right: value{
			val: delta,
		},
	}
}
func (c Column) Multi(delta int) MathExpr {
	return MathExpr{
		left: c,
		opt:  optMULTI,
		right: value{
			val: delta,
		},
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

// InQuery 一种是 IN 子查询, 另外一种就是普通的值
func (c Column) InQuery(sub Subquery) Predicate {
	return Predicate{
		left:  c,
		opt:   optIN,
		right: sub,
	}
}
