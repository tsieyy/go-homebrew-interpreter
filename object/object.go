package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ = "ERROR_OBJ"  // TODO：这里是不是有问题
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ = "STRING"
	BUILTIN_OBJ = "BUILTIN"
	ARRAY_OBJ = "ARRAY"
	HASH_OBJ = "HASH"
	QUOTE_OBJ = "QUOTE"
)


// 数值
type Integer struct {
	Value int64
}
//
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}


// bool
type Boolean struct {
	Value bool
}
//
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}


// null
type Null struct {}
//
func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}


// return
type ReturnValue struct {
	Value Object
}
//
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}


// error
type Error struct {
	Message string
}
//
func (e *Error) Inspect() string {
	return "Error:" + e.Message
}
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}


// function
type Function struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Environment
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn") 
	out.WriteString("(") 
	out.WriteString(strings.Join(params, ", ")) 
	out.WriteString(") {\n") 
	out.WriteString(f.Body.String()) 
	out.WriteString("\n}") 
	return out.String()
}
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}


type String struct {
	Value string
}
func (s *String) Inspect() string {
	return s.Value
}
func (s *String) Type() ObjectType {
	return STRING_OBJ
}



type BuiltinFunction func (args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}
func (b *Builtin) Inspect() string {
	return "built-in function"
}
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}


type Array struct {
	Elements []Object
}
func (a *Array) Inspect() string {
	var out bytes.Buffer 
	elements := []string{} 
	for _, e := range a.Elements { 
		elements = append(elements, e.Inspect()) 
	} 
	out.WriteString("[") 
	out.WriteString(strings.Join(elements, ", ")) 
	out.WriteString("]")
	return out.String()
}
func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}


type Hashable interface {
	HashKey() HashKey
}


type HashKey struct {
	Type ObjectType
	Value uint64
}
// 给三个基本的数据结构加上hashKey方法
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{
		Type: b.Type(),
		Value: value,
	}
}
func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type: i.Type(),
		Value: uint64(i.Value),
	}
}
// TODO:该方法存在哈希碰撞的问题，可以尝试解决
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{
		Type: s.Type(),
		Value: h.Sum64(),
	}
}


type HashPair struct {
	Key Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}
func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}
func (h *Hash) Inspect() string {
	var out bytes.Buffer 
	pairs := []string{} 
	for _, pair := range h.Pairs { 
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect())) 
	} 
	out.WriteString("{") 
	out.WriteString(strings.Join(pairs, ", ")) 
	out.WriteString("}") 
	return out.String()
}


type Quote struct {
	Node ast.Node
}
func (q *Quote) Type() ObjectType {
	return QUOTE_OBJ
}
func (q *Quote) Inspect() string {
	return "Quote(" + q.Node.String() + ")"
}
