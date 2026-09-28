//go:build js && wasm

// Package main provides the WebAssembly (WASM) entry point for MOJA interpreter.
// It exports the runMojoLang JavaScript function for executing MOJA code in browser.
package main

import (
	"bytes"
	"syscall/js"

	"mojolang/internal/evaluator"
	"mojolang/internal/lexer"
	"mojolang/internal/parser"
	"mojolang/internal/reporter"
)

func runMojoLang(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"output":  "",
			"error":   "[MOJA Error] Kode sumber tidak diberikan",
			"success": false,
		}
	}

	source := args[0].String()

	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		errStr := err.Error()
		if posErr, ok := err.(reporter.PositionalError); ok {
			line, col := posErr.GetPosition()
			errStr = reporter.FormatError("playground.mojo", source, line, col, posErr.GetErrorType(), err.Error())
		}
		return map[string]interface{}{
			"output":  "",
			"error":   errStr,
			"success": false,
		}
	}

	p := parser.New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		errStr := err.Error()
		if posErr, ok := err.(reporter.PositionalError); ok {
			line, col := posErr.GetPosition()
			errStr = reporter.FormatError("playground.mojo", source, line, col, posErr.GetErrorType(), err.Error())
		}
		return map[string]interface{}{
			"output":  "",
			"error":   errStr,
			"success": false,
		}
	}

	var outBuf bytes.Buffer
	ev := evaluator.New(&outBuf)
	env := evaluator.NewEnvironment()
	res := ev.Eval(program, env)

	if errObj, ok := res.(*evaluator.RuntimeError); ok {
		line, col := errObj.GetPosition()
		errStr := reporter.FormatError("playground.mojo", source, line, col, errObj.GetErrorType(), errObj.Error())
		return map[string]interface{}{
			"output":  outBuf.String(),
			"error":   errStr,
			"success": false,
		}
	}

	return map[string]interface{}{
		"output":  outBuf.String(),
		"error":   "",
		"success": true,
	}
}

func main() {
	js.Global().Set("runMojoLang", js.FuncOf(runMojoLang))
	// Keep Go runtime alive in browser
	select {}
}
