// Package repl implements the Read Eval Print Loop function
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/actuallyachraf/monkey-giggle/eval"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/parser"

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
	env := object.NewEnv()
	fmt.Print(WELCOME)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Printf(EXIT)
			break
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.Parse()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaled := eval.Eval(program, env)
		if evaled != nil {
			io.WriteString(out, evaled.Inspect())
			io.WriteString(out, "\n")
		}

	}

}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
