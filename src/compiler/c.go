package compiler

import (
	"bytes"
	"os"
	"os/exec"
)

func GetCCode(source []byte) bytes.Buffer {
	var buff bytes.Buffer
	var programStart = `
	extern int putchar(int);
	extern char getchar();

	char memory[30000];
	int ptr = 0;

	int main () {
	`
	buff.WriteString(programStart)

	for _, bt := range source {
		switch bt {
		case 'i':
			buff.WriteString("ptr++;\n")
		case 'a':
			buff.WriteString("ptr--;\n")
		case 'l':
			buff.WriteString("memory[ptr]++;\n")
		case 'e':
			buff.WriteString("memory[ptr]--;\n")
		case 'o':
			buff.WriteString("putchar(memory[ptr]);\n")
		case 'v':
			buff.WriteString("memory[ptr] = getchar();\n")
		case 'b':
			buff.WriteString("while (memory[ptr]) {\n")
		case 'u':
			buff.WriteString("}\n")
		}
	}
	buff.WriteString("return 0;}")

	return buff
}

func writeToFile(buffer *bytes.Buffer) (string, error) {
	err := os.WriteFile("a.c", buffer.Bytes(), 0644)
	return "a.c", err
}

func GenerateCCode(source []byte) (string, error) {
	buff := GetCCode(source)
	filename, err := writeToFile(&buff)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func CompileCCode(filename string) (string, error) {
	cmd := exec.Command("gcc", "-o", "a.out", filename)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "a.out", nil
}

func RunCCode(filename string) error {
	cmd := exec.Command("./a.out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
