package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ASC8384/wen/src/compiler"
)

func main() {

	compile := flag.Bool("c", false, "compile")
	flag.Parse()

	// read file
	args := os.Args
	if len(args) == 2 {
		filename := args[1]
		code, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", filename)
			return
		}
		// execute
		compiler.Execute(string(code), filename)
		return
	}
	if len(args) != 3 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}
	filename := args[2]
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filename)
		return
	}
	if *compile {
		var cfilename string
		cfilename, err = compiler.GenerateCCode(code)
		compiler.CompileCCode(cfilename)
		compiler.RunCCode(cfilename)
		// ccode := GetCCode(code)
		// cfilename, err := writeToFile(&ccode)
		// if err != nil {
		// 	os.Exit(1)
		// }
		// output := "./a.out"
		// gcc := exec.Command("gcc", "-O3", "-Ofast", "-o", output, cfilename)
		// gcc.Stdout = os.Stdout
		// gcc.Stderr = os.Stderr
		// err = gcc.Run()
		// if err != nil {
		// 	fmt.Printf("Error launching gcc: %s\n", err)
		// 	return
		// }
		// exe := exec.Command(output)
		// exe.Stdout = os.Stdout
		// exe.Stderr = os.Stderr
		// err = exe.Run()
		// if err != nil {
		// 	fmt.Printf("Error launching %s: %s\n", output, err)
		// 	os.Exit(1)
		// }
	}

}
