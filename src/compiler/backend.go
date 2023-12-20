package compiler

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/parser"
)

const MEM_SIZE int = 30000

type MemoryVariable struct {
	Start  int
	Length int
}

type GlobalVariables struct {
	Variables map[string]MemoryVariable
}

type GlobalMemory struct {
	Memory []int
}

func NewGlobalVariables() *GlobalVariables {
	var g GlobalVariables
	g.Variables = make(map[string]MemoryVariable)
	return &g
}

func NewGlobalMemory() *GlobalMemory {
	var m GlobalMemory
	m.Memory = make([]int, MEM_SIZE)
	return &m
}

func NewGlobalPointer() *int {
	var ptr int
	return &ptr
}

func Execute(code, filename string) {
	var ast *Block
	var err error

	g := NewGlobalVariables()
	m := NewGlobalMemory()
	p := NewGlobalPointer()

	// parse
	if ast, err = Parse(code, filename); err != nil {
		panic(err)
	}

	// resolve
	if err = resolveAST(g, m, p, ast); err != nil {
		panic(err)
	}
}

func resolveAST(g *GlobalVariables, m *GlobalMemory, p *int, ast *Block) error {
	if len(ast.Stats) == 0 {
		return errors.New("resolveAST(): no code to execute, please check your input.")
	}
	// Execute init statements first
	// for _, statement := range ast.Stats {
	switch s := ast.Stats[len(ast.Stats)-1].(type) {
	case *StringExp:
		if err := resolveInit(m, s); err != nil {
			return err
		}
		break
	}
	// }
	// Execute all other statements
	for _, statement := range ast.Stats {
		if err := resolveStatement(g, m, p, statement); err != nil {
			return err
		}
	}
	return nil
}

func resolveStatement(g *GlobalVariables, m *GlobalMemory, p *int, statement Stat) error {
	// fmt.Printf("statement: %T\n", statement)
	switch s := statement.(type) {
	case *AssignStat:
		return resolveAssignment(g, m, p, s)
	case *Print:
		return resolvePrint(g, m, p, s)
	case *Scanf:
		return resolveScanf(g, m, p, s)
	case *PointerStat:
		return resolvePointer(p, s)
	case *CellStat:
		return resolveCell(m, p, s)
	case *LoopStat:
		return resolveLoop(g, m, p, s)
	case *StringExp:
		return nil
	default:
		return fmt.Errorf("resolveStatement(): undefined statement type: %T", statement)
	}
}

func resolveInit(m *GlobalMemory, init *StringExp) error {
	Length := len(init.Str)
	for i := 0; i < Length; i++ {
		m.Memory[i] = int(init.Str[i])
	}
	return nil
}

func resolveLoop(g *GlobalVariables, m *GlobalMemory, p *int, loop *LoopStat) error {
	for m.Memory[*p] != 0 { // Continue looping while the current cell is not zero
		for _, stat := range loop.Stats {
			err := resolveStatement(g, m, p, stat) // Resolve each statement in the loop body
			if err != nil {
				return err // Return if an error occurs
			}

			// It's important to check bounds of the memory and pointer
			if *p < 0 || *p >= len(m.Memory) {
				return errors.New("resolveLoop(): data pointer out of bounds")
			}
		}
	}
	return nil
}

func resolveAssignment(g *GlobalVariables, m *GlobalMemory, p *int, assignment *AssignStat) error {
	varName := ""
	if varName = assignment.Variable.Name; varName == "" {
		return errors.New("resolveAssignment(): variable name can NOT be empty.")
	}
	Length := len(assignment.String)
	if Length != 0 {
		if assignment.Int {
			val, err := strconv.Atoi(assignment.String)
			if err != nil {
				return err
			}
			m.Memory[*p] = val
			g.Variables[varName] = MemoryVariable{Start: *p, Length: 1}
			*p += 1
		} else {
			for i := 0; i < Length; i++ {
				m.Memory[*p+i] = int(assignment.String[i])
			}
			g.Variables[varName] = MemoryVariable{Start: *p, Length: Length}
			*p += Length
		}
	} else {
		if assignment.Length != 0 {
			// Length = assignment.Length
			var gVar = g.Variables[varName]
			gVar.Length += assignment.Length
			g.Variables[varName] = gVar
		} else if assignment.Start != 0 {
			var gVar = g.Variables[varName]
			gVar.Start += assignment.Start
			g.Variables[varName] = gVar
		}
	}
	return nil
}

func resolvePrint(g *GlobalVariables, m *GlobalMemory, p *int, print *Print) error {
	varName := ""
	if nil == print.Variable {
		if print.Int {
			fmt.Printf("%d", m.Memory[*p])
		} else {
			fmt.Printf("%c", m.Memory[*p])
		}
		return nil
	} else if varName = print.Variable.Name; varName == "" {
		return errors.New("resolvePrint(): variable name can NOT be empty.")
	} else {
		Start := g.Variables[varName].Start
		Length := g.Variables[varName].Length
		if print.Int {
			for i := 0; i < Length; i++ {
				fmt.Printf("%d", m.Memory[Start+i])
			}
		} else {
			for i := 0; i < Length; i++ {
				fmt.Printf("%c", m.Memory[Start+i])
			}
		}
	}
	return nil
}

func resolveScanf(g *GlobalVariables, m *GlobalMemory, p *int, scanf *Scanf) error {
	varName := ""
	if nil == scanf.Variable {
		if scanf.Int {
			var input int
			fmt.Scanf("%d", &input)
			m.Memory[*p] = input
		} else {
			buf := make([]byte, 1)
			len, err := os.Stdin.Read(buf)
			if err != nil {
				return err
			}
			if len != 1 {
				return fmt.Errorf("read %d bytes of input, not 1", len)
			}
			m.Memory[*p] = int(buf[0])
		}
		return nil
	} else if varName = scanf.Variable.Name; varName == "" {
		return errors.New("resolvePrint(): variable name can NOT be empty.")
	} else {
		if scanf.Int {
			var input int
			fmt.Scanf("%d", &input)
			m.Memory[*p] = input
			g.Variables[varName] = MemoryVariable{Start: *p, Length: 1}
			*p += 1
		} else {
			var input string
			fmt.Scanf("%s", &input)
			Length := len(input)
			for i := 0; i < Length; i++ {
				m.Memory[*p+i] = int(input[i])
			}
			g.Variables[varName] = MemoryVariable{Start: *p, Length: Length}
			*p += Length
		}
		// g.Variables[varName] = string(input)
	}
	return nil
}

func resolvePointer(p *int, pointer *PointerStat) error {
	*p += pointer.Pointer
	return nil
}

func resolveCell(m *GlobalMemory, p *int, cell *CellStat) error {
	m.Memory[*p] += cell.Cell
	return nil
}
