// Package ast defines the Abstract Syntax Tree (AST) nodes for MOJA.
package ast

import (
	"bytes"
	"fmt"
	"strings"
)

// Node interface for all AST nodes.
type Node interface {
	String() string
	LineNum() int
	ColNum() int
}

// Statement interface for statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression interface for expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every MOJA AST.
type Program struct {
	Statements []Statement
}

func (p *Program) LineNum() int {
	if len(p.Statements) > 0 {
		return p.Statements[0].LineNum()
	}
	return 1
}

func (p *Program) ColNum() int {
	if len(p.Statements) > 0 {
		return p.Statements[0].ColNum()
	}
	return 1
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

// Identifier node
type Identifier struct {
	Line   int
	Column int
	Value  string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) LineNum() int    { return i.Line }
func (i *Identifier) ColNum() int     { return i.Column }
func (i *Identifier) String() string  { return i.Value }

// NumberLiteral node
type NumberLiteral struct {
	Line   int
	Column int
	Value  float64
	Raw    string
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) LineNum() int    { return n.Line }
func (n *NumberLiteral) ColNum() int     { return n.Column }
func (n *NumberLiteral) String() string  { return n.Raw }

// StringLiteral node
type StringLiteral struct {
	Line   int
	Column int
	Value  string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) LineNum() int    { return s.Line }
func (s *StringLiteral) ColNum() int     { return s.Column }
func (s *StringLiteral) String() string  { return fmt.Sprintf("%q", s.Value) }

// BooleanLiteral node
type BooleanLiteral struct {
	Line   int
	Column int
	Value  bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) LineNum() int    { return b.Line }
func (b *BooleanLiteral) ColNum() int     { return b.Column }
func (b *BooleanLiteral) String() string {
	if b.Value {
		return "bener"
	}
	return "salah"
}

// PrefixExpression (-x)
type PrefixExpression struct {
	Line     int
	Column   int
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) LineNum() int    { return pe.Line }
func (pe *PrefixExpression) ColNum() int     { return pe.Column }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// InfixExpression (x + y, x == y, etc.)
type InfixExpression struct {
	Line     int
	Column   int
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) LineNum() int    { return ie.Line }
func (ie *InfixExpression) ColNum() int     { return ie.Column }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// CallExpression (f(a, b))
type CallExpression struct {
	Line      int
	Column    int
	Function  *Identifier // Function name
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) LineNum() int    { return ce.Line }
func (ce *CallExpression) ColNum() int     { return ce.Column }
func (ce *CallExpression) String() string {
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}

// OmongStatement: omong(expr)
type OmongStatement struct {
	Line   int
	Column int
	Value  Expression
}

func (os *OmongStatement) statementNode() {}
func (os *OmongStatement) LineNum() int   { return os.Line }
func (os *OmongStatement) ColNum() int    { return os.Column }
func (os *OmongStatement) String() string { return fmt.Sprintf("omong(%s)", os.Value.String()) }

// AssignStatement: x = expr
type AssignStatement struct {
	Line   int
	Column int
	Name   *Identifier
	Value  Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) LineNum() int   { return as.Line }
func (as *AssignStatement) ColNum() int    { return as.Column }
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Name.String(), as.Value.String())
}

// BlockStatement: { statements }
type BlockStatement struct {
	Line       int
	Column     int
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) LineNum() int   { return bs.Line }
func (bs *BlockStatement) ColNum() int    { return bs.Column }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// YenStatement: yen condition block (liyane block)?
type YenStatement struct {
	Line        int
	Column      int
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // can be nil
}

func (ys *YenStatement) statementNode() {}
func (ys *YenStatement) LineNum() int   { return ys.Line }
func (ys *YenStatement) ColNum() int    { return ys.Column }
func (ys *YenStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("yen %s %s", ys.Condition.String(), ys.Consequence.String()))
	if ys.Alternative != nil {
		out.WriteString(fmt.Sprintf(" liyane %s", ys.Alternative.String()))
	}
	return out.String()
}

// BaleniStatement: baleni var saka start tekan end block
type BaleniStatement struct {
	Line     int
	Column   int
	Iterator *Identifier
	Start    Expression
	End      Expression
	Body     *BlockStatement
}

func (bs *BaleniStatement) statementNode() {}
func (bs *BaleniStatement) LineNum() int   { return bs.Line }
func (bs *BaleniStatement) ColNum() int    { return bs.Column }
func (bs *BaleniStatement) String() string {
	return fmt.Sprintf("baleni %s saka %s tekan %s %s", bs.Iterator.String(), bs.Start.String(), bs.End.String(), bs.Body.String())
}

// GaweStatement: gawe name(params) block
type GaweStatement struct {
	Line   int
	Column int
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStatement
}

func (gs *GaweStatement) statementNode() {}
func (gs *GaweStatement) LineNum() int   { return gs.Line }
func (gs *GaweStatement) ColNum() int    { return gs.Column }
func (gs *GaweStatement) String() string {
	var params []string
	for _, p := range gs.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("gawe %s(%s) %s", gs.Name.String(), strings.Join(params, ", "), gs.Body.String())
}

// ReturnStatement: balik expr
type ReturnStatement struct {
	Line   int
	Column int
	Value  Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) LineNum() int   { return rs.Line }
func (rs *ReturnStatement) ColNum() int    { return rs.Column }
func (rs *ReturnStatement) String() string { return fmt.Sprintf("balik %s", rs.Value.String()) }

// ExpressionStatement: standalone expression
type ExpressionStatement struct {
	Line       int
	Column     int
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) LineNum() int   { return es.Line }
func (es *ExpressionStatement) ColNum() int    { return es.Column }
func (es *ExpressionStatement) String() string { return es.Expression.String() }

// ArrayLiteral: [elem1, elem2, ...]
type ArrayLiteral struct {
	Line     int
	Column   int
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) LineNum() int    { return al.Line }
func (al *ArrayLiteral) ColNum() int     { return al.Column }
func (al *ArrayLiteral) String() string {
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// IndexExpression: arr[index]
type IndexExpression struct {
	Line   int
	Column int
	Left   Expression
	Index  Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) LineNum() int    { return ie.Line }
func (ie *IndexExpression) ColNum() int     { return ie.Column }
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}

// IndexAssignStatement: arr[index] = value
type IndexAssignStatement struct {
	Line   int
	Column int
	Left   *IndexExpression
	Value  Expression
}

func (ias *IndexAssignStatement) statementNode() {}
func (ias *IndexAssignStatement) LineNum() int   { return ias.Line }
func (ias *IndexAssignStatement) ColNum() int    { return ias.Column }
func (ias *IndexAssignStatement) String() string {
	return fmt.Sprintf("%s = %s", ias.Left.String(), ias.Value.String())
}
