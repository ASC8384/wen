package compiler

import (
	"errors"
	"fmt"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/parser"
)

type GlobalVariables struct {
	Variables map[string]string
}

func NewGlobalVariables() *GlobalVariables {
	var g GlobalVariables
	g.Variables = make(map[string]string)
	return &g
}

func Execute(code, filename string) {
	var ast *Block
	var err error

	g := NewGlobalVariables()

	// parse
	if ast, err = Parse(code, filename); err != nil {
		panic(err)
	}

	// resolve
	if err = resolveAST(g, ast); err != nil {
		panic(err)
	}
}

func resolveAST(g *GlobalVariables, ast *Block) error {
	if len(ast.Stats) == 0 {
		return errors.New("resolveAST(): no code to execute, please check your input.")
	}
	for _, statement := range ast.Stats {
		if err := resolveStatement(g, statement); err != nil {
			return err
		}
	}
	return nil
}

func resolveStatement(g *GlobalVariables, statement Stat) error {
	if assignment, ok := statement.(*AssignStat); ok {
		return resolveAssignment(g, assignment)
	} else if print, ok := statement.(*Print); ok {
		return resolvePrint(g, print)
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

func resolvePrint(g *GlobalVariables, print *Print) error {
	varName := ""
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
