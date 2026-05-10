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

type Stmt interface {
	stmtNode()
	Location() gears.Location
	String() string
	Accept(Visitor) (any, error)
}

type Expr interface {
	exprNode()
	Location() gears.Location
	String() string
	Accept(Visitor) (any, error)
}

type Type interface {
	typeNode()
	Location() gears.Location
	String() string
	Accept(Visitor) (any, error)
}
