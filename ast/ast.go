package ast

import (
	"bytes"
	"monkey/token"
)

// 该接口用于调试
type Node interface {
	// 返回与其关联的词法单元的字面量
	TokenLiteral() string
	// 调试的时候打印 AST信息
	String() string 
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
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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
func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString("=")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}




type ReturnStatement struct {
	Token token.Token //token.RETURN
	ReturnValue Expression
}
// 实现接口
func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}


// 除去let、return 就是表达式语句，例如：x+10；这种，在一些脚本语言中比较常见了
type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}
// 实现接口
func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) String() string {
	//var out bytes.Buffer
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
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
func (i *Identifier) String() string {
	return i.Value
}



type IntegerLiteral struct { 
	Token token.Token 
	Value int64 
} 
// 实现接口
func (il *IntegerLiteral) expressionNode() {} 
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal } 
func (il *IntegerLiteral) String() string { return il.Token.Literal }



type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}
// 实现接口
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}


type InfixExpression struct {
	Token token.Token
	Operator string
	Left Expression
	Right Expression 
}
// 实现接口
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}


