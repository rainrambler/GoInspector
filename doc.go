package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var curfset *token.FileSet

func scanSrcDir(dirname string) {
	err := filepath.Walk(dirname, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.go", filepath.Base(fpath)); err != nil {
			return err
		} else if matched {
			ParseGoSrc(fpath)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ParseGoSrc(fname string) {
	fmt.Printf("INFO: Parsing %s...\n", fname)
	bs, err := ReadBinFile(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	fset := token.NewFileSet()
	curfset = fset

	f, err := parser.ParseFile(fset, "main.go", string(bs), parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	//printDecl(f, fset)

	var v visitor
	ast.Walk(v, f)
}

type visitor struct{}

// https://adlerhsieh.com/blog/write-your-own-go-linters-with-parser-package
// Visit function walks through each node in a file
func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	blockStmt, ok := n.(*ast.BlockStmt)
	if ok {
		parseBlockStmt(blockStmt)
	}

	return v
}

func parseBlockStmt(stmt *ast.BlockStmt) {
	if len(stmt.List) < 4 {
		// do not need to parse
		return
	}

	var nps NamePairs

	for _, child := range stmt.List {
		switch nodetp := child.(type) {
		case *ast.AssignStmt:
			np := convNamePair(nodetp)
			nps.Add(np)
		default:
		}
	}

	nps.findExchange()
}

type NamePair struct {
	NPos  token.Pos
	Left  string
	Right string
}

func convNamePair(assStmt *ast.AssignStmt) *NamePair {
	if len(assStmt.Lhs) != 1 {
		// complex left variant
		return nil
	}
	lval, ok := assStmt.Lhs[0].(*ast.Ident)
	if !ok {
		fmt.Printf("WARN: Cannot convert left [%+v]!\n", assStmt)
		return nil
	}

	np := new(NamePair)
	np.Left = lval.Name
	np.NPos = lval.NamePos

	if len(assStmt.Rhs) != 1 {
		// complex right variant
		return nil
	}
	rval, ok := assStmt.Rhs[0].(*ast.Ident)
	if !ok {
		//fmt.Printf("WARN: Cannot convert right [%v]!\n",
		//	curfset.Position(assStmt.Pos()))
		return nil
	}
	np.Right = rval.Name
	return np
}

type NamePairs struct {
	arr []*NamePair
}

func (p *NamePairs) Add(np *NamePair) {
	if np == nil {
		return
	}
	p.arr = append(p.arr, np)
}

func (p *NamePairs) findExchange() {
	if len(p.arr) < 3 {
		return
	}

	total := len(p.arr)

	for i := 0; i < total-2; i++ {
		first := p.arr[i]
		second := p.arr[i+1]
		third := p.arr[i+2]

		if (first.Right == second.Left) &&
			(second.Right == third.Left) &&
			(third.Right == first.Left) {
			fmt.Printf("Find exchange in [%v]!\n", curfset.Position(first.NPos))
		}
	}
}
