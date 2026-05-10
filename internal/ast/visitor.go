package ast

import "go/ast"

type Visitor interface {
	Visit(files []ast.File) (any, error)
	visitExpr(Expr) (any, error)
	visitAtomExpr(*AtomExpr) (any, error)
	visitUnaryExpr(*UnaryExpr) (any, error)
	visitBinaryExpr(*BinaryExpr) (any, error)
	visitFunctionCallExpr(*FunctionCallExpr) (any, error)
	visitStmt(Stmt) (any, error)
	visitExprStmt(*ExprStmt) (any, error)
	visitPrintStmt(*PrintStmt) (any, error)
	visitStructDefStmt(*StructDefStmt) (any, error)
	visitReturnStmt(*ReturnStmt) (any, error)
	visitFuncDefStmt(*FuncDefStmt) (any, error)
	visitVarDefStmt(*VarDefStmt) (any, error)
	visitModuleDefStmt(*ModuleDefStmt) (any, error)
	visitType(Type) (any, error)
	visitNamedT(*NamedT) (any, error)
	visitIntT(*IntT) (any, error)
	visitFloatT(*FloatT) (any, error)
	visitStringT(*StringT) (any, error)
}
