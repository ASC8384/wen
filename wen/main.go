package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ASC8384/wen/src/compiler"
)

func main() {

	compile := flag.Bool("c", false, "compile")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()

	// read file
	args := os.Args
	var filename string
	if !*compile {
		if *debug {
			filename = args[2]
		} else {
			filename = args[1]
		}
		code, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", filename)
			return
		}
		// execute
		compiler.Execute(string(code), filename, *debug)
		return
	}
	if len(args) != 3 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}
	filename = args[2]
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filename)
		return
	}
	if *compile {
		var cfilename string
		cfilename, err = compiler.GenerateCCode(code)
		if err != nil {
			fmt.Printf("Error generating C code: %s\n", filename)
			return
		}
		compiler.CompileCCode(cfilename)
		compiler.RunCCode(cfilename)
	}

}
