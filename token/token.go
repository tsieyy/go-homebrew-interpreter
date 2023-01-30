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
	ASSIGN = "=" 
	PLUS = "+" 
	MINUS = "-" 
	BANG = "!" 
	ASTERISK = "*" 
	SLASH = "/"

	LT = "<" 
	GT = ">"

	EQ = "=="
	NOT_EQ = "!="

	// 分隔符
	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "(" 
	RPAREN = ")" 
	LBRACE = "{" 
	RBRACE = "}"
	LBRACKET = "[" 
 	RBRACKET = "]"

	// 关键字
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
	STRING = "STRING"

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


// 关键字的映射，用于判断字符是关键字还是标识符
var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

// 查找是否在keyword中，以判断是否是关键字还是标识符
func LookupIdent(ident string) TokenType {
	// t , ok := keywords[ident]
	if t, ok := keywords[ident]; ok {
		return t
	} else {
		return IDENT
	}
}

