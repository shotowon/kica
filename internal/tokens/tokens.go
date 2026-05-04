package tokens

import (
	"fmt"
	"slices"

	"github.com/shotowon/kica/internal/gears"
)

type TokenKind int

func (k TokenKind) IsOneOf(kinds ...TokenKind) bool {
	return slices.Contains(kinds, k)
}

const (
	EOF TokenKind = iota
	Module
	Struct
	Function
	Return
	Var
	Print
	TLInt
	TLFloat
	TLString
	LParen // (
	RParen // )
	LBrace // {
	RBrace // }

	Plus
	Minus
	Mul
	Div
	Eq
	LT // <
	GT // >
	Colon
	Semicolon
	Comma
	ID
	String
	Int
	Float
)

type Token struct {
	Location gears.Location
	Kind     TokenKind
	Literal  string
}

func (t Token) String() string {
	return fmt.Sprintf("<Token(%s): literal: '%s' at %s>", t.Kind, t.Literal, t.Location)
}

func (t TokenKind) String() string {
	switch t {
	case EOF:
		return "EOF"
	case Module:
		return "module"
	case Function:
		return "func"
	case Return:
		return "return"
	case Struct:
		return "struct"
	case Print:
		return "print"
	case LParen:
		return "("
	case RParen:
		return ")"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Mul:
		return "*"
	case Div:
		return "/"
	case ID:
		return "ID"
	case String:
		return "string"
	case Int:
		return "int"
	case Float:
		return "float"
	case TLString:
		return "tl-string"
	case TLInt:
		return "tl-int"
	case TLFloat:
		return "tl-float"
	case Colon:
		return ":"
	case Semicolon:
		return ";"
	case GT:
		return ">"
	case LT:
		return "<"
	case Comma:
		return ","
	default:
		return fmt.Sprintf("TokenKind(%d)", int(t))
	}
}

func IDtoKeyword(literal string) TokenKind {
	switch literal {
	case "function":
		return Function
	case "struct":
		return Struct
	case "return":
		return Return
	case "module":
		return Module
	case "print":
		return Print
	case "var":
		return Var
	case "string":
		return TLString
	case "float":
		return TLFloat
	case "int":
		return TLInt
	}

	return ID
}
