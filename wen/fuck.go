package main

import (
	"fmt"
	"os"
)

func interpret(code string) {
	memory := make([]byte, 30000)
	ptr := 0

	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '>':
			ptr++
		case '<':
			ptr--
		case '+':
			memory[ptr]++
		case '-':
			memory[ptr]--
		case '.':
			fmt.Print(string(memory[ptr]))
		case ',':
			var input byte
			fmt.Scanf("%c", &input)
			memory[ptr] = input
		case '[':
			if memory[ptr] == 0 {
				loopCount := 1
				for loopCount > 0 {
					i++
					if code[i] == '[' {
						loopCount++
					} else if code[i] == ']' {
						loopCount--
					}
				}
			}
		case ']':
			loopCount := 1
			for loopCount > 0 {
				i--
				if code[i] == '[' {
					loopCount--
				} else if code[i] == ']' {
					loopCount++
				}
			}
			i--
		}
	}
}

func main1() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <brainfuck_code>")
		return
	}

	code := os.Args[1]
	interpret(code)
}
