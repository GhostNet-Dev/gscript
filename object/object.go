package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/GhostNet-Dev/glambda/ast"
	"github.com/GhostNet-Dev/glambda/code"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	IDENTFIER_OBJ         = "IDENTIFIER"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	STRING_OBJ            = "STRING"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	HASH_OBJ              = "HASH"
	CLOSURE_OBJ           = "CLOSURE"
)

type Hashable interface{ HashKey() HashKey }

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (o *Hash) Type() ObjectType { return HASH_OBJ }
func (o *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range o.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Array struct {
	Elements []Object
}

func (o *Array) Type() ObjectType { return ARRAY_OBJ }
func (o *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range o.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (o *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (o *Builtin) Inspect() string  { return "builtin function" }

type String struct {
	Value string
}

func (o *String) Type() ObjectType { return STRING_OBJ }
func (o *String) Inspect() string  { return o.Value }
func (o *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(o.Value))
	return HashKey{Type: o.Type(), Value: h.Sum64()}
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (o *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (o *Closure) Inspect() string  { return fmt.Sprintf("Closure[%p]", o) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (o *Function) Type() ObjectType { return FUNCTION_OBJ }
func (o *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range o.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(o.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (o *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }
func (o *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", o)
}

type Error struct {
	Message string
}

func (o *Error) Inspect() string  { return "ERROR: " + o.Message }
func (o *Error) Type() ObjectType { return ERROR_OBJ }

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }
func (o *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Identifier struct {
	Name  string
	Value Object
}

func (o *Identifier) Inspect() string  { return fmt.Sprintf("%d", o.Value) }
func (o *Identifier) Type() ObjectType { return IDENTFIER_OBJ }

type Integer struct {
	Value int64
}

func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }
func (o *Integer) Type() ObjectType { return INTEGER_OBJ }
func (o *Integer) HashKey() HashKey { return HashKey{Type: o.Type(), Value: uint64(o.Value)} }

type Boolean struct {
	Value bool
}

func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }
func (o *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (o *Boolean) HashKey() HashKey {
	var value uint64
	if o.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: o.Type(), Value: value}
}

type Null struct{}

func (o *Null) Inspect() string  { return "null" }
func (o *Null) Type() ObjectType { return NULL_OBJ }
