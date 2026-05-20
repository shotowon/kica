package ast

import "github.com/shotowon/kica/internal/gears"

type File struct {
	stmts []Stmt
}

func NewFile(stmts []Stmt) *File {
	return &File{
		stmts: stmts,
	}
}

func (f *File) Statements() []Stmt {
	return f.stmts
}

type Node interface {
	Location() gears.Location
	String() string
	Accept(Visitor) (any, error)
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Type interface {
	Node
	typeNode()
}
