package ast

import "github.com/shotowon/kica/internal/gears"

type Stmt interface {
	stmtNode()
	Location() gears.Location
	String() string
}

type Expr interface {
	exprNode()
	Location() gears.Location
	String() string
}

type Type interface {
	typeNode()
	Location() gears.Location
	String() string
}
