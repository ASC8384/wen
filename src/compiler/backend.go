package compiler

import (
	"errors"
	"fmt"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/parser"
)

const MEM_SIZE int = 30000

type GlobalVariables struct {
	Variables map[string]string
}

type GlobalMemory struct {
	Memory []int
}

func NewGlobalVariables() *GlobalVariables {
	var g GlobalVariables
	g.Variables = make(map[string]string)
	return &g
}

func NewGlobalMemory() *GlobalMemory {
	var m GlobalMemory
	m.Memory = make([]int, MEM_SIZE)
	m.Memory[2] = 2222
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
	for _, statement := range ast.Stats {
		if err := resolveStatement(g, m, p, statement); err != nil {
			return err
		}
	}
	return nil
}

func resolveStatement(g *GlobalVariables, m *GlobalMemory, p *int, statement Stat) error {
	if assignment, ok := statement.(*AssignStat); ok {
		return resolveAssignment(g, assignment)
	} else if print, ok := statement.(*Print); ok {
		return resolvePrint(g, m, p, print)
	} else if pointer, ok := statement.(*PointerStat); ok {
		return resolvePointer(p, pointer)
	} else if cell, ok := statement.(*CellStat); ok {
		return resolveCell(m, p, cell)
	} else {
		return errors.New("resolveStatement(): undefined statement type.")
	}
}

func resolveAssignment(g *GlobalVariables, assignment *AssignStat) error {
	varName := ""
	if varName = assignment.Variable.Name; varName == "" {
		return errors.New("resolveAssignment(): variable name can NOT be empty.")
	}
	g.Variables[varName] = assignment.String
	return nil
}

func resolvePrint(g *GlobalVariables, m *GlobalMemory, p *int, print *Print) error {
	varName := ""
	if nil == print.Variable {
		fmt.Print(m.Memory[*p])
		return nil
	}
	if varName = print.Variable.Name; varName == "" {
		return errors.New("resolvePrint(): variable name can NOT be empty.")
	}
	str := ""
	ok := false
	if str, ok = g.Variables[varName]; !ok {
		return errors.New(fmt.Sprintf("resolvePrint(): variable '$%s'not found.", varName))
	}
	fmt.Print(str)
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
