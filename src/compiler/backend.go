package compiler

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	. "github.com/ASC8384/wen/src/ast"
	"github.com/ASC8384/wen/src/lexer"
	. "github.com/ASC8384/wen/src/parser"
)

const MEM_SIZE int = 30000

// MemoryVariable represents a memory variable with a start index and length.
type MemoryVariable struct {
	Start  int
	Length int
}

// GlobalVariables represents the global variables in the program.
type GlobalVariables struct {
	Variables map[string]MemoryVariable
}

// GlobalMemory represents the global memory in the program.
type GlobalMemory struct {
	Memory []int
}

// NewGlobalVariables creates a new instance of GlobalVariables.
func NewGlobalVariables() *GlobalVariables {
	var g GlobalVariables
	g.Variables = make(map[string]MemoryVariable)
	return &g
}

// NewGlobalMemory creates a new instance of GlobalMemory.
func NewGlobalMemory() *GlobalMemory {
	var m GlobalMemory
	m.Memory = make([]int, MEM_SIZE)
	return &m
}

// NewGlobalPointer creates a new instance of a global pointer.
func NewGlobalPointer() *int {
	var ptr int
	return &ptr
}

// NewGlobalRegister creates a new instance of a global register.
func NewGlobalRegister() *int {
	var reg int
	return &reg
}

func printAST(node []Stat, indent string) {
	if node == nil {
		return
	}

	fmt.Println(indent, node)

	for _, stat := range node {
		switch sta := stat.(type) {
		case *Print:
			fmt.Println(indent, "Print")
			fmt.Println(indent+"  ", sta.Variable)
			fmt.Println(indent+"  ", sta.Int)
		case *Scanf:
			fmt.Println(indent, "Scanf")
			fmt.Println(indent+"  ", sta.Variable)
			fmt.Println(indent+"  ", sta.Int)
		case *AssignStat:
			fmt.Println(indent, "AssignStat")
			fmt.Println(indent+"  ", sta.Variable)
			fmt.Println(indent+"  ", sta.String)
			fmt.Println(indent+"  ", sta.Int)
			fmt.Println(indent+"  ", sta.Length)
		case *PointerStat:
			fmt.Println(indent, "PointerStat")
			fmt.Println(indent+"  ", sta.Pointer)
		case *CellStat:
			fmt.Println(indent, "CellStat")
			fmt.Println(indent+"  ", sta.Cell)
		case *LoopStat:
			fmt.Println(indent, "LoopStat")
			printAST(sta.Stats, indent+"|    ")
		case *RegStat:
			fmt.Println(indent, "RegStat")
			fmt.Println(indent+"  ", sta.Reg)
		case *IfStat:
			fmt.Println(indent, "IfStat")
			fmt.Println(indent+"  ", sta.Type)
		}
	}
}

// Execute executes the given code with the specified filename.
func Execute(code, filename string, debug bool) {
	var ast *Block
	var err error

	g := NewGlobalVariables()
	m := NewGlobalMemory()
	p := NewGlobalPointer()
	r := NewGlobalRegister()

	// parse
	if ast, err = Parse(code, filename); err != nil {
		panic(err)
	}

	if debug {
		// print Abstract Syntax Tree
		fmt.Println("AST:")
		// fmt.Printf("%#v\n", ast)
		printAST(ast.Stats, "|")
	}

	// resolve
	if err = resolveAST(g, m, p, r, ast); err != nil {
		panic(err)
	}
}

// resolveAST resolves the abstract syntax tree (AST) of the program.
func resolveAST(g *GlobalVariables, m *GlobalMemory, p *int, r *int, ast *Block) error {
	if len(ast.Stats) == 0 {
		return errors.New("resolveAST(): no code to execute, please check your input")
	}
	// Execute init statements first
	switch s := ast.Stats[len(ast.Stats)-1].(type) {
	case *StringExp:
		if err := resolveInit(m, s); err != nil {
			return err
		}
	}
	// Execute all other statements
	// continueCnt is a variable that represents the number of times a loop should continue.
	var continueCnt int
	// rushB is a boolean that represents whether we are in a rush B status.
	var rushB bool
	for _, statement := range ast.Stats {
		if continueCnt > 0 {
			continueCnt--
			continue
		}
		if err := resolveStatement(rushB, g, m, p, r, statement); err != nil {
			if err.Error() == ("resolveIf(): rushB") { // If we encounter a rush B statement, we should skip the next statement.
				rushB = true
				continue
			} else if err.Error() == ("resolveIf(): end rush B") && rushB { // If we encounter an end rush B statement, we should stop skipping the next statement.
				rushB = false
				continueCnt = 0
				continue
			} else if err.Error() == ("resolveIf(): time to die") {
				return nil
			} else if err.Error() == "resolveIf(): start == 0" {
				continueCnt = 0
				continue
			} else if err.Error() == ("resolveIf(): start > 0") {
				continueCnt = 1
				continue
			} else if err.Error() == ("resolveIf(): start < 0") {
				continueCnt = 2
				continue
			}
			return err
		}
	}
	return nil
}

// resolveStatement resolves a statement in the program.
func resolveStatement(rushB bool, g *GlobalVariables, m *GlobalMemory, p *int, r *int, statement Stat) error {
	if rushB {
		switch statement := statement.(type) {
		case *IfStat:
			return resolveIf(m, p, statement)
		default:
			return nil
		}
	}
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
		return resolveLoop(rushB, g, m, p, r, s)
	case *RegStat:
		return resolveReg(m, p, r, s)
	case *IfStat:
		return resolveIf(m, p, s)
	case *StringExp:
		return nil
	default:
		return fmt.Errorf("resolveStatement(): undefined statement type: %T", statement)
	}
}

// resolveIf resolves an if statement in the program.
func resolveIf(m *GlobalMemory, p *int, ifStat *IfStat) error {
	if ifStat.Type == lexer.TOKEN_IF_START {
		if m.Memory[*p] == 0 { // If the current cell is zero, we should skip the next statement.
			*p += 1
			return errors.New("resolveIf(): start == 0")
		} else if m.Memory[*p] > 0 { // If the current cell is greater than zero, we should skip the next two statements.
			*p += 2
			return errors.New("resolveIf(): start > 0")
		} else { // If the current cell is less than zero, we should skip the next three statements.
			*p += 3
			return errors.New("resolveIf(): start < 0")
		}
	} else if ifStat.Type == lexer.TOKEN_IF_RUSH {
		return errors.New("resolveIf(): rushB")
	} else if ifStat.Type == lexer.TOKEN_IF_END {
		return errors.New("resolveIf(): end rush B")
	} else if ifStat.Type == lexer.TOKEN_END {
		return errors.New("resolveIf(): time to die")
	}
	return nil
}

// resolveReg resolves a register statement in the program.
func resolveReg(m *GlobalMemory, p *int, r *int, reg *RegStat) error {
	switch reg.Reg {
	case lexer.TOKEN_REG_STORE:
		// store
		*r = m.Memory[*p]
	case lexer.TOKEN_REG_PLUS:
		// plus
		m.Memory[*p] += *r
	case lexer.TOKEN_REG_MINUS:
		// minus
		m.Memory[*p] -= *r
	case lexer.TOKEN_REG_MUL:
		// mul
		m.Memory[*p] *= *r
	case lexer.TOKEN_REG_DIV:
		// div
		m.Memory[*p] /= *r
	case lexer.TOKEN_REG_MOD:
		// mod
		m.Memory[*p] %= *r
	case lexer.TOKEN_REG_READ:
		// read
		m.Memory[*p] = *r
	default:
		return fmt.Errorf("resolveReg(): undefined reg type in %T: %d", reg, reg.Reg)
	}
	return nil
}

// resolveInit resolves an initialization statement in the program.
func resolveInit(m *GlobalMemory, init *StringExp) error {
	Length := len(init.Str)
	for i := 0; i < Length; i++ {
		m.Memory[i] = int(init.Str[i])
	}
	return nil
}

// resolveLoop resolves a loop statement in the program.
func resolveLoop(rushB bool, g *GlobalVariables, m *GlobalMemory, p *int, r *int, loop *LoopStat) error {
	for m.Memory[*p] != 0 { // Continue looping while the current cell is not zero
		for _, stat := range loop.Stats {
			err := resolveStatement(rushB, g, m, p, r, stat) // Resolve each statement in the loop body
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

// resolveAssignment resolves an assignment statement in the program.
func resolveAssignment(g *GlobalVariables, m *GlobalMemory, p *int, assignment *AssignStat) error {
	varName := ""
	if varName = assignment.Variable.Name; varName == "" {
		return errors.New("resolveAssignment(): variable name can NOT be empty")
	}
	Length := len(assignment.String)
	if Length != 0 { // if assignment.String is not empty
		if assignment.Int { // using int format
			val, err := strconv.Atoi(assignment.String)
			if err != nil {
				return err
			}
			m.Memory[*p] = val
			g.Variables[varName] = MemoryVariable{Start: *p, Length: 1}
			*p += 1
		} else { // using string format
			for i := 0; i < Length; i++ {
				m.Memory[*p+i] = int(assignment.String[i])
			}
			g.Variables[varName] = MemoryVariable{Start: *p, Length: Length}
			*p += Length
		}
	} else { // adjust the length or start of the variable
		if assignment.Length != 0 { // adjust the length of the variable
			var gVar = g.Variables[varName]
			gVar.Length += assignment.Length
			g.Variables[varName] = gVar
		} else if assignment.Start != 0 { // adjust the start of the variable
			var gVar = g.Variables[varName]
			gVar.Start += assignment.Start
			g.Variables[varName] = gVar
		} else { // get the start of the variable
			// TOKEN_VAR_START
			*p = g.Variables[varName].Start
		}
	}
	return nil
}

// resolvePrint resolves a print statement in the program.
func resolvePrint(g *GlobalVariables, m *GlobalMemory, p *int, print *Print) error {
	varName := ""
	if nil == print.Variable { // print value at the data pointer
		if print.Int {
			fmt.Printf("%d", m.Memory[*p])
		} else {
			fmt.Printf("%c", m.Memory[*p])
		}
		return nil
	} else if varName = print.Variable.Name; varName == "" {
		return errors.New("resolvePrint(): variable name can NOT be empty")
	} else { // print value of the variable
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

// resolveScanf resolves a scanf statement in the program.
func resolveScanf(g *GlobalVariables, m *GlobalMemory, p *int, scanf *Scanf) error {
	varName := ""
	if nil == scanf.Variable { // scanf value at the data pointer
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
		return errors.New("resolvePrint(): variable name can NOT be empty")
	} else { // scanf value of the variable
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

// resolvePointer resolves a pointer statement in the program.
func resolvePointer(p *int, pointer *PointerStat) error {
	*p += pointer.Pointer
	return nil
}

// resolveCell resolves a cell statement in the program.
func resolveCell(m *GlobalMemory, p *int, cell *CellStat) error {
	m.Memory[*p] += cell.Cell
	return nil
}
