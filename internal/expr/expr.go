package expr

type Expr interface {
	isExpr()
}
type Number int

func (i Number) isExpr() {}

type Label string

func (i Label) isExpr() {}

type BinOp struct {
	Left  Expr
	Op    string
	Right Expr
}

func (i BinOp) isExpr() {}
