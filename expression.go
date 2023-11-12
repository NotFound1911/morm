package morm

// RawExpr 原生表达式，不做任何处理
type RawExpr struct {
	raw  string
	args []interface{}
}

func (e RawExpr) selectable() {
}

func (e RawExpr) expr() {

}

func (e RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: e,
	}
}

func Raw(expr string, args ...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}
