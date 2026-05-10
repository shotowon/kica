package ast

import (
	"fmt"
	"strings"

	"github.com/shotowon/kica/internal/gears"
	"github.com/shotowon/kica/internal/tokens"
)

type AtomExpr struct {
	Value *tokens.Token
}

func (e *AtomExpr) exprNode() {}
func (e *AtomExpr) Location() gears.Location {
	return e.Value.Location
}

func (e *AtomExpr) String() string {
	if e.Value == nil {
		return "<nil>"
	}
	return e.Value.Literal
}

func (s *AtomExpr) Accept(v Visitor) (any, error) {
	return v.visitAtomExpr(s)
}

type UnaryExpr struct {
	Value Expr
	Op    *tokens.Token
}

func (e *UnaryExpr) exprNode() {}
func (e *UnaryExpr) Location() gears.Location {
	loc := e.Op.Location
	loc.End = e.Value.Location().End
	return loc
}

func (e *UnaryExpr) String() string {
	op := "?"
	if e.Op != nil {
		op = e.Op.Literal
	}
	return fmt.Sprintf("(%s %v)", op, e.Value)
}

func (s *UnaryExpr) Accept(v Visitor) (any, error) {
	return v.visitUnaryExpr(s)
}

type BinaryExpr struct {
	LHS Expr
	RHS Expr
	Op  *tokens.Token
}

func (e *BinaryExpr) exprNode() {}
func (e *BinaryExpr) Location() gears.Location {
	loc := e.LHS.Location()
	loc.End = e.RHS.Location().End
	return loc
}

func (e *BinaryExpr) String() string {
	op := "?"
	if e.Op != nil {
		op = e.Op.Literal
	}
	return fmt.Sprintf("(%s %v %v)", op, e.LHS, e.RHS)
}

func (s *BinaryExpr) Accept(v Visitor) (any, error) {
	return v.visitBinaryExpr(s)
}

type FunctionCallExpr struct {
	Name *tokens.Token
	Args []Expr
	Loc  gears.Location
}

func (e *FunctionCallExpr) exprNode() {}
func (e *FunctionCallExpr) Location() gears.Location {
	return e.Loc
}

func (e *FunctionCallExpr) String() string {
	args := strings.Builder{}
	for i, arg := range e.Args {
		args.WriteString(arg.String())
		if i+1 < len(e.Args) {
			args.WriteString(", ")
		}
	}
	return fmt.Sprintf("(%s(%s))", e.Name.Literal, args.String())
}

func (s *FunctionCallExpr) Accept(v Visitor) (any, error) {
	return v.visitFunctionCallExpr(s)
}
