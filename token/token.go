// 词法分析过程中需用到的 词法单元（Token）



package token



const (
	EOF = "EOF"
	ILLEGAL = "ILLEGAL"

	// 标识符
	IDENT = "IDENT"  //变量名等等

	// 字面量
	INT = "INT" // 123  321

	// 运算符
	PLUS = "+"
	ASSIGN = "="

	// 分隔符
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "(" 
	RPAREN = ")" 
	LBRACE = "{" 
	RBRACE = "}"

	// 关键字
	FUNCTION = "FUNCTION"
	LET = "LET"

)



type TokenType string


type Token struct {
	Type TokenType
	Literal string
}


func NewToken(t TokenType, b byte) *Token {
	return &Token{
		Type: t,
		Literal: string(b),
	}
}