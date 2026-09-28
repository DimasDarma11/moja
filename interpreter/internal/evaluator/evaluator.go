// Package evaluator implements the runtime evaluation for MOJA AST.
package evaluator

import (
	"fmt"
	"io"
	"strconv"

	"mojolang/internal/ast"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Evaluator struct {
	out        io.Writer
	inFunction bool
	callDepth  int
}

func New(out io.Writer) *Evaluator {
	return &Evaluator{
		out:        out,
		inFunction: false,
		callDepth:  0,
	}
}

// Eval evaluates an AST node in the given environment.
func (ev *Evaluator) Eval(node ast.Node, env *Environment) Object {
	if node == nil {
		return nil
	}

	switch node := node.(type) {

	case *ast.Program:
		return ev.evalProgram(node, env)

	case *ast.BlockStatement:
		return ev.evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return ev.Eval(node.Expression, env)

	case *ast.OmongStatement:
		val := ev.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		fmt.Fprintln(ev.out, val.Inspect())
		return nil

	case *ast.AssignStatement:
		val := ev.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val, node.Line, node.Column)

	case *ast.IndexAssignStatement:
		return ev.evalIndexAssignStatement(node, env)

	case *ast.YenStatement:
		return ev.evalYenStatement(node, env)

	case *ast.BaleniStatement:
		return ev.evalBaleniStatement(node, env)

	case *ast.GaweStatement:
		fn := &Function{
			Name:   node.Name.Value,
			Params: node.Params,
			Body:   node.Body,
			Env:    env,
		}
		env.Set(node.Name.Value, fn, node.Line, node.Column)
		return nil

	case *ast.ReturnStatement:
		if !ev.inFunction {
			return &RuntimeError{
				Line:    node.Line,
				Column:  node.Column,
				Message: "'balik' hanya dapat digunakan di dalam fungsi",
			}
		}
		val := ev.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *ast.Identifier:
		return ev.evalIdentifier(node, env)

	case *ast.NumberLiteral:
		return &Number{Value: node.Value}

	case *ast.StringLiteral:
		return &String{Value: node.Value}

	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.ArrayLiteral:
		return ev.evalArrayLiteral(node, env)

	case *ast.IndexExpression:
		left := ev.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := ev.Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return ev.evalIndexExpression(left, index, node.Line, node.Column)

	case *ast.PrefixExpression:
		right := ev.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return ev.evalPrefixExpression(node.Operator, right, node.Line, node.Column)

	case *ast.InfixExpression:
		// Short-circuit untuk operator logika: evaluasi left terlebih dulu
		if node.Operator == "lan" || node.Operator == "utawa" {
			return ev.evalLogicalInfix(node, env)
		}
		left := ev.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := ev.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return ev.evalInfixExpression(node.Operator, left, right, node.Line, node.Column)

	case *ast.CallExpression:
		return ev.evalCallExpression(node, env)
	}

	return nil
}

func (ev *Evaluator) evalProgram(program *ast.Program, env *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = ev.Eval(statement, env)

		switch res := result.(type) {
		case *ReturnValue:
			return &RuntimeError{
				Line:    statement.LineNum(),
				Column:  statement.ColNum(),
				Message: "'balik' hanya dapat digunakan di dalam fungsi",
			}
		case *RuntimeError:
			return res
		}
	}
	return result
}

func (ev *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object
	for _, statement := range block.Statements {
		result = ev.Eval(statement, env)
		if result != nil {
			if result.Type() == RETURN_VALUE || result.Type() == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func (ev *Evaluator) evalIdentifier(node *ast.Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("Variabel '%s' belum dideklarasikan", node.Value),
		}
	}
	return val
}

func (ev *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *Environment) Object {
	var elements []Object
	for _, elNode := range node.Elements {
		elVal := ev.Eval(elNode, env)
		if isError(elVal) {
			return elVal
		}
		elements = append(elements, elVal)
	}
	return &Array{Elements: elements}
}

func (ev *Evaluator) evalIndexExpression(leftObj, indexObj Object, line, col int) Object {
	arr, ok := leftObj.(*Array)
	if !ok {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: fmt.Sprintf("Akses indeks '[...]' hanya berlaku untuk tipe Array, diberikan %s", leftObj.Type()),
		}
	}

	num, ok := indexObj.(*Number)
	if !ok {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: fmt.Sprintf("Indeks array harus berupa tipe Number, diberikan %s", indexObj.Type()),
		}
	}

	if num.Value != float64(int(num.Value)) {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: "Indeks array harus berupa bilangan bulat (integer)",
		}
	}

	idx := int(num.Value)
	if idx < 0 {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: "Indeks array tidak boleh negatif",
		}
	}

	if idx >= len(arr.Elements) {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: fmt.Sprintf("Indeks array di luar rentang (posisi %d, total elemen %d)", idx, len(arr.Elements)),
		}
	}

	return arr.Elements[idx]
}

func (ev *Evaluator) evalIndexAssignStatement(node *ast.IndexAssignStatement, env *Environment) Object {
	leftObj := ev.Eval(node.Left.Left, env)
	if isError(leftObj) {
		return leftObj
	}

	indexObj := ev.Eval(node.Left.Index, env)
	if isError(indexObj) {
		return indexObj
	}

	val := ev.Eval(node.Value, env)
	if isError(val) {
		return val
	}

	arr, ok := leftObj.(*Array)
	if !ok {
		return &RuntimeError{
			Line:    node.Left.Line,
			Column:  node.Left.Column,
			Message: fmt.Sprintf("Mutasi indeks '[...]' hanya berlaku untuk tipe Array, diberikan %s", leftObj.Type()),
		}
	}

	num, ok := indexObj.(*Number)
	if !ok {
		return &RuntimeError{
			Line:    node.Left.Line,
			Column:  node.Left.Column,
			Message: fmt.Sprintf("Indeks array harus berupa tipe Number, diberikan %s", indexObj.Type()),
		}
	}

	if num.Value != float64(int(num.Value)) {
		return &RuntimeError{
			Line:    node.Left.Line,
			Column:  node.Left.Column,
			Message: "Indeks array harus berupa bilangan bulat (integer)",
		}
	}

	idx := int(num.Value)
	if idx < 0 {
		return &RuntimeError{
			Line:    node.Left.Line,
			Column:  node.Left.Column,
			Message: "Indeks array tidak boleh negatif",
		}
	}

	if idx >= len(arr.Elements) {
		return &RuntimeError{
			Line:    node.Left.Line,
			Column:  node.Left.Column,
			Message: fmt.Sprintf("Indeks array di luar rentang (posisi %d, total elemen %d)", idx, len(arr.Elements)),
		}
	}

	arr.Elements[idx] = val
	return nil
}

func (ev *Evaluator) evalPrefixExpression(op string, right Object, line, col int) Object {
	switch op {
	case "-":
		num, ok := right.(*Number)
		if !ok {
			return &RuntimeError{
				Line:    line,
				Column:  col,
				Message: "Operator '-' hanya mendukung tipe Number",
			}
		}
		return &Number{Value: -num.Value}
	case "ora":
		bool, ok := right.(*Boolean)
		if !ok {
			return &RuntimeError{
				Line:    line,
				Column:  col,
				Message: fmt.Sprintf("Operator 'ora' hanya mendukung tipe Boolean (bener/salah), diberikan %s", right.Type()),
			}
		}
		return nativeBoolToBooleanObject(!bool.Value)
	}
	return &RuntimeError{Line: line, Column: col, Message: fmt.Sprintf("Operator prefix tidak dikenal: %s", op)}
}

func (ev *Evaluator) evalInfixExpression(op string, left, right Object, line, col int) Object {
	// String concatenation rules for +
	if op == "+" {
		if left.Type() == STRING_OBJ || right.Type() == STRING_OBJ {
			return &String{Value: left.Inspect() + right.Inspect()}
		}
		leftNum, okL := left.(*Number)
		rightNum, okR := right.(*Number)
		if okL && okR {
			return &Number{Value: leftNum.Value + rightNum.Value}
		}
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: "Operator '+' membutuhkan tipe Number atau String",
		}
	}

	// Equality comparisons work for all types
	if op == "==" {
		return nativeBoolToBooleanObject(ev.evalEquals(left, right))
	}
	if op == "!=" {
		return nativeBoolToBooleanObject(!ev.evalEquals(left, right))
	}

	// Arithmetic (-, *, /) and Relational (>, <, >=, <=) require Number
	leftNum, okL := left.(*Number)
	rightNum, okR := right.(*Number)
	if !okL || !okR {
		return &RuntimeError{
			Line:    line,
			Column:  col,
			Message: fmt.Sprintf("Operator '%s' hanya mendukung tipe Number", op),
		}
	}

	switch op {
	case "-":
		return &Number{Value: leftNum.Value - rightNum.Value}
	case "*":
		return &Number{Value: leftNum.Value * rightNum.Value}
	case "/":
		if rightNum.Value == 0 {
			return &RuntimeError{
				Line:    line,
				Column:  col,
				Message: "Pembagian dengan nol tidak diperbolehkan",
			}
		}
		return &Number{Value: leftNum.Value / rightNum.Value}
	case ">":
		return nativeBoolToBooleanObject(leftNum.Value > rightNum.Value)
	case "<":
		return nativeBoolToBooleanObject(leftNum.Value < rightNum.Value)
	case ">=":
		return nativeBoolToBooleanObject(leftNum.Value >= rightNum.Value)
	case "<=":
		return nativeBoolToBooleanObject(leftNum.Value <= rightNum.Value)
	}

	return &RuntimeError{Line: line, Column: col, Message: fmt.Sprintf("Operator tidak dikenal: %s", op)}
}

// evalLogicalInfix mengevaluasi operator 'lan' dan 'utawa' dengan short-circuit.
// 'lan': jika kiri salah → kembalikan salah tanpa mengevaluasi kanan.
// 'utawa': jika kiri bener → kembalikan bener tanpa mengevaluasi kanan.
func (ev *Evaluator) evalLogicalInfix(node *ast.InfixExpression, env *Environment) Object {
	left := ev.Eval(node.Left, env)
	if isError(left) {
		return left
	}
	leftBool, ok := left.(*Boolean)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("Operand kiri operator '%s' harus bertipe Boolean (bener/salah), diberikan %s", node.Operator, left.Type()),
		}
	}

	// Short-circuit evaluation
	if node.Operator == "lan" && !leftBool.Value {
		return FALSE // salah lan ... → salah (tanpa evaluasi kanan)
	}
	if node.Operator == "utawa" && leftBool.Value {
		return TRUE // bener utawa ... → bener (tanpa evaluasi kanan)
	}

	right := ev.Eval(node.Right, env)
	if isError(right) {
		return right
	}
	rightBool, ok := right.(*Boolean)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("Operand kanan operator '%s' harus bertipe Boolean (bener/salah), diberikan %s", node.Operator, right.Type()),
		}
	}

	switch node.Operator {
	case "lan":
		return nativeBoolToBooleanObject(leftBool.Value && rightBool.Value)
	case "utawa":
		return nativeBoolToBooleanObject(leftBool.Value || rightBool.Value)
	}

	return &RuntimeError{Line: node.Line, Column: node.Column, Message: fmt.Sprintf("Operator logika tidak dikenal: %s", node.Operator)}
}

func (ev *Evaluator) evalEquals(left, right Object) bool {
	if left.Type() != right.Type() {
		return false
	}
	switch l := left.(type) {
	case *Number:
		r := right.(*Number)
		return l.Value == r.Value
	case *String:
		r := right.(*String)
		return l.Value == r.Value
	case *Boolean:
		r := right.(*Boolean)
		return l.Value == r.Value
	}
	return false
}

func (ev *Evaluator) evalYenStatement(node *ast.YenStatement, env *Environment) Object {
	cond := ev.Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	condBool, ok := cond.(*Boolean)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: "Kondisi 'yen' harus bernilai Boolean (bener/salah)",
		}
	}

	blockEnv := NewEnclosedEnvironment(env)
	if condBool.Value {
		return ev.Eval(node.Consequence, blockEnv)
	} else if node.Alternative != nil {
		return ev.Eval(node.Alternative, blockEnv)
	}

	return nil
}

func (ev *Evaluator) evalBaleniStatement(node *ast.BaleniStatement, env *Environment) Object {
	startObj := ev.Eval(node.Start, env)
	if isError(startObj) {
		return startObj
	}
	startNum, ok := startObj.(*Number)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: "Batas awal perulangan 'baleni' harus berupa Number",
		}
	}

	endObj := ev.Eval(node.End, env)
	if isError(endObj) {
		return endObj
	}
	endNum, ok := endObj.(*Number)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: "Batas akhir perulangan 'baleni' harus berupa Number",
		}
	}

	startVal := int(startNum.Value)
	endVal := int(endNum.Value)

	if startVal > endVal {
		return nil
	}

	for i := startVal; i <= endVal; i++ {
		loopEnv := NewEnclosedEnvironment(env)
		loopEnv.DefineIterator(node.Iterator.Value, &Number{Value: float64(i)})

		result := ev.Eval(node.Body, loopEnv)
		if result != nil {
			if result.Type() == RETURN_VALUE || result.Type() == ERROR_OBJ {
				return result
			}
		}
	}

	return nil
}

func (ev *Evaluator) evalCallExpression(node *ast.CallExpression, env *Environment) Object {
	fnObj, ok := env.Get(node.Function.Value)
	if !ok {
		// Builtin function 'panjang'
		if node.Function.Value == "panjang" {
			if len(node.Arguments) != 1 {
				return &RuntimeError{
					Line:    node.Line,
					Column:  node.Column,
					Message: fmt.Sprintf("Fungsi 'panjang' butuh 1 argumen, diberikan %d", len(node.Arguments)),
				}
			}
			argVal := ev.Eval(node.Arguments[0], env)
			if isError(argVal) {
				return argVal
			}
			if arr, ok := argVal.(*Array); ok {
				return &Number{Value: float64(len(arr.Elements))}
			}
			if str, ok := argVal.(*String); ok {
				return &Number{Value: float64(len([]rune(str.Value)))}
			}
			return &RuntimeError{
				Line:    node.Line,
				Column:  node.Column,
				Message: fmt.Sprintf("Fungsi 'panjang' hanya mendukung tipe Array atau String, diberikan %s", argVal.Type()),
			}
		}

		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("Fungsi '%s' belum dideklarasikan", node.Function.Value),
		}
	}

	fn, ok := fnObj.(*Function)
	if !ok {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("'%s' bukan merupakan fungsi", node.Function.Value),
		}
	}

	if len(node.Arguments) != len(fn.Params) {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: fmt.Sprintf("Fungsi '%s' butuh %d argumen, diberikan %d", fn.Name, len(fn.Params), len(node.Arguments)),
		}
	}

	if ev.callDepth >= 1000 {
		return &RuntimeError{
			Line:    node.Line,
			Column:  node.Column,
			Message: "Kedalaman rekursi melampaui batas (1000)",
		}
	}

	var args []Object
	for _, argNode := range node.Arguments {
		argVal := ev.Eval(argNode, env)
		if isError(argVal) {
			return argVal
		}
		args = append(args, argVal)
	}

	// Prepare function execution environment with closure env
	fnEnv := NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Params {
		fnEnv.Set(param.Value, args[i], node.Line, node.Column)
	}

	prevInFunction := ev.inFunction
	ev.inFunction = true
	ev.callDepth++

	result := ev.Eval(fn.Body, fnEnv)

	ev.callDepth--
	ev.inFunction = prevInFunction

	if rv, ok := result.(*ReturnValue); ok {
		return rv.Value
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
