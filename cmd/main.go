package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/actuallyachraf/monkey-giggle/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s !\n", user.Username)
	fmt.Printf("Feel free to type in commands !\n")
	repl.Start(os.Stdin, os.Stdout)
}
