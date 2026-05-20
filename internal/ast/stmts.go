package ast

import (
	"fmt"
	"strings"

	"github.com/shotowon/kica/internal/gears"
	"github.com/shotowon/kica/internal/tokens"
)

type ExprStmt struct {
	Value Expr
	Loc   gears.Location
}

func (s *ExprStmt) stmtNode() {}
func (s *ExprStmt) Location() gears.Location {
	return s.Loc
}

func (s *ExprStmt) String() string {
	if s.Value == nil {
		return "(expr <nil>)"
	}
	return fmt.Sprintf("(expr %s)", s.Value)
}

func (s *ExprStmt) Accept(v Visitor) (any, error) {
	return v.VisitExprStmt(s)
}

type PrintStmt struct {
	Loc   gears.Location
	Value Expr
}

func (s *PrintStmt) stmtNode() {}
func (s *PrintStmt) Location() gears.Location {
	return s.Loc
}

func (s *PrintStmt) String() string {
	if s.Value == nil {
		return "(print <nil>)"
	}
	return fmt.Sprintf("(print %s)", s.Value)
}

func (s *PrintStmt) Accept(v Visitor) (any, error) {
	return v.VisitPrintStmt(s)
}

type StructDefStmt struct {
	Loc    gears.Location
	Name   *tokens.Token
	Fields map[string]Type
}

func (s *StructDefStmt) stmtNode() {}
func (s *StructDefStmt) Location() gears.Location {
	return s.Loc
}

func (s *StructDefStmt) String() string {
	var fieldsStr strings.Builder
	for fieldName, fieldType := range s.Fields {
		fmt.Fprintf(&fieldsStr, "%s: %s\n", fieldName, fieldType)
	}

	return fmt.Sprintf("{struct %s %s}", s.Name, fieldsStr.String())
}

func (s *StructDefStmt) Accept(v Visitor) (any, error) {
	return v.VisitStructDefStmt(s)
}

type FuncDefStmt struct {
	Loc        gears.Location
	Name       *tokens.Token
	Params     []NTMapping
	ReturnType Type
	Body       []Stmt
}

func (s *FuncDefStmt) stmtNode() {}
func (s *FuncDefStmt) Location() gears.Location {
	return s.Loc
}

func (s *FuncDefStmt) String() string {
	name := "<nil>"
	if s.Name != nil {
		name = s.Name.Literal
	}

	return fmt.Sprintf("function %s (%s) -> %s\n{\n%s\n}",
		name,
		gears.Map(s.Params, func(nt NTMapping) string {
			return fmt.Sprintf("%s: %s ", nt.Name, nt.Ty)
		}),
		s.ReturnType,
		gears.Map(s.Body, func(s Stmt) string { return fmt.Sprintf("%s ", s) }))
}

func (s *FuncDefStmt) Accept(v Visitor) (any, error) {
	return v.VisitFuncDefStmt(s)
}

type ReturnStmt struct {
	Value Expr
	Loc   gears.Location
}

func (s *ReturnStmt) stmtNode() {}
func (s *ReturnStmt) Location() gears.Location {
	return s.Loc
}

func (s *ReturnStmt) String() string {
	if s.Value != nil {
		return fmt.Sprintf("return %s", s.Value)
	}

	return "return"
}

func (s *ReturnStmt) Accept(v Visitor) (any, error) {
	return v.VisitReturnStmt(s)
}

type VarDefStmt struct {
	Name  *tokens.Token
	Value Expr
	Loc   gears.Location
}

func (s *VarDefStmt) stmtNode() {}
func (s *VarDefStmt) Location() gears.Location {
	return s.Loc
}

func (s *VarDefStmt) String() string {
	if s.Value == nil {
		return fmt.Sprintf("(var %s <nil>)", s.Name.Literal)
	}
	return fmt.Sprintf("(var %s %s)", s.Name.Literal, s.Value)
}

func (s *VarDefStmt) Accept(v Visitor) (any, error) {
	return v.VisitVarDefStmt(s)
}

type ModuleDefStmt struct {
	Loc  gears.Location
	Name *tokens.Token
}

func (s *ModuleDefStmt) stmtNode() {}
func (s *ModuleDefStmt) Location() gears.Location {
	return s.Loc
}

func (s *ModuleDefStmt) String() string {
	return fmt.Sprintf("package %s", s.Name.Literal)
}

func (s *ModuleDefStmt) Accept(v Visitor) (any, error) {
	return v.VisitModuleDefStmt(s)
}
