// Package repl implements the Read Eval Print Loop function
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/actuallyachraf/monkey-giggle/compiler"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/parser"
	"github.com/actuallyachraf/monkey-giggle/vm"

	"github.com/actuallyachraf/monkey-giggle/lexer"
)

// PROMPT marks prompt level console
const PROMPT = "giggle>> "

// WELCOME the user to make us giggle
const WELCOME = "Make me giggle !\n"

// EXIT the repl
const EXIT = "Ohhh you're leaving already !"

// Start the read eval print loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltIn(i, v.Name)
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println(EXIT)
			break
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalState(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
