package morm

// Expr 原生表达式，不做任何处理
type Expr struct {
	raw  string
	args []interface{}
}

func (e Expr) selectable() {
}

func (e Expr) expr() {

}

func (e Expr) AsPredicate() Predicate {
	return Predicate{
		left: e,
	}
}

func Raw(expr string, args ...interface{}) Expr {
	return Expr{
		raw:  expr,
		args: args,
	}
}
