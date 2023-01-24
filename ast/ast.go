package ast

import (
	"monkey/token"
)


// 该接口用于调试
type Node interface {
	// 返回与其关联的词法单元的字面量
	TokenLiteral() string
}

// 语句
type Statement interface {
	Node
	statementNode() 
}

// 表达式
type Expression interface {
	Node
	expressionNode() 
}



// 存语句
type Program struct {
	Statements []Statement
}
// 实现接口
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}


type LetStatement struct {
	Token token.Token //token.LET
	Name *Identifier
	Value Expression
}
// 实现接口
func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}



type Identifier struct {
	Token token.Token //token.IDENT
	Value string
}
// 实现接口，属于表达式了,因为有些情况下，标识符会产生值，比如把一个值赋值给另一个标识符的时候
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
} 