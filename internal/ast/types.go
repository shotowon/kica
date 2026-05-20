package ast

import (
	"fmt"

	"github.com/shotowon/kica/internal/gears"
	"github.com/shotowon/kica/internal/tokens"
)

type NamedT struct {
	Tok *tokens.Token
}

func (t *NamedT) typeNode() {}
func (t *NamedT) Location() gears.Location {
	return t.Tok.Location
}

func (t *NamedT) String() string {
	return t.Tok.Literal
}

func (t *NamedT) Accept(v Visitor) (any, error) {
	return v.VisitNamedT(t)
}

type IntT struct {
	Loc gears.Location
}

func (t *IntT) typeNode() {}
func (t *IntT) Location() gears.Location {
	return t.Loc
}

func (t *IntT) String() string {
	return "int"
}

func (t *IntT) Accept(v Visitor) (any, error) {
	return v.VisitIntT(t)
}

type FloatT struct {
	Loc gears.Location
}

func (t *FloatT) typeNode() {}
func (t *FloatT) Location() gears.Location {
	return t.Loc
}

func (t *FloatT) String() string {
	return "float"
}

func (t *FloatT) Accept(v Visitor) (any, error) {
	return v.VisitFloatT(t)
}

type StringT struct {
	Loc gears.Location
}

func (t *StringT) typeNode() {}
func (t *StringT) Location() gears.Location {
	return t.Loc
}

func (t *StringT) String() string {
	return "string"
}

func (t *StringT) Accept(v Visitor) (any, error) {
	return v.VisitStringT(t)
}

type FunctionT struct {
	Loc        gears.Location
	ParamsType []Type
	ReturnType Type
}

func (t *FunctionT) typeNode() {}
func (t *FunctionT) Location() gears.Location {
	return t.Loc
}

func (t *FunctionT) String() string {
	return fmt.Sprintf(
		"function (%s) -> %s",
		gears.Map(t.ParamsType, func(pt Type) string { return pt.String() }),
		t.ReturnType.String(),
	)
}

func (t *FunctionT) Accept(v Visitor) (any, error) {
	return v.VisitFunctionT(t)
}

type StructT struct {
	Loc    gears.Location
	Name   tokens.Token
	Fields map[string]Type
}

func (t *StructT) typeNode() {}
func (t *StructT) Location() gears.Location {
	return t.Loc
}

func (t *StructT) String() string {
	return "struct %s"
}

func (t *StructT) Accept(v Visitor) (any, error) {
	return v.VisitStructT(t)
}
