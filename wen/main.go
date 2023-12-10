package main

import (
	"fmt"
	"os"

	"github.com/ASC8384/wen/src/compiler"
)

func main() {

	// read file
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}
	filename := args[1]
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filename)
		return
	}

	// execute
	compiler.Execute(string(code), filename)

}
