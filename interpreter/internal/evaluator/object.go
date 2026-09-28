package evaluator

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"mojolang/internal/ast"
)

type ObjectType string

const (
	NUMBER_OBJ   ObjectType = "NUMBER"
	STRING_OBJ   ObjectType = "STRING"
	BOOLEAN_OBJ  ObjectType = "BOOLEAN"
	ARRAY_OBJ    ObjectType = "ARRAY"
	RETURN_VALUE ObjectType = "RETURN_VALUE"
	FUNCTION_OBJ ObjectType = "FUNCTION"
	ERROR_OBJ    ObjectType = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Number object
type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return strconv.FormatFloat(n.Value, 'f', -1, 64) }

// String object
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "bener"
	}
	return "salah"
}

// Array object
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var elements []string
	for _, e := range a.Elements {
		if e != nil {
			if str, ok := e.(*String); ok {
				elements = append(elements, fmt.Sprintf("%q", str.Value))
			} else {
				elements = append(elements, e.Inspect())
			}
		} else {
			elements = append(elements, "nil")
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// ReturnValue wrapper object
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Function object
type Function struct {
	Name   string
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Params {
		params = append(params, p.Value)
	}
	var out bytes.Buffer
	out.WriteString("gawe ")
	out.WriteString(f.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

// RuntimeError object
type RuntimeError struct {
	Line    int
	Column  int
	Message string
}

func (e *RuntimeError) Type() ObjectType { return ERROR_OBJ }
func (e *RuntimeError) Inspect() string  { return e.Error() }
func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[MOJA Error] Baris %d, Kolom %d: %s", e.Line, e.Column, e.Message)
}

func (e *RuntimeError) GetPosition() (int, int) {
	return e.Line, e.Column
}

func (e *RuntimeError) GetErrorType() string {
	return "Runtime Error"
}
