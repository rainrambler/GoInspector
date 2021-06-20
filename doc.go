package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var curfset *token.FileSet

func ParseGoSrc(fname string) {
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

	// Print the AST.
	ast.Print(fset, f) // https://zhuanlan.zhihu.com/p/28516587
	var v visitor
	ast.Walk(v, f)

	//递归调用逐一打印节点
	/*
		// https://www.jianshu.com/p/443bd82863f8
		ast.Inspect(f, func(n ast.Node) bool {
			ast.Print(fset, n)
			return true
		})
	*/

}

type visitor struct{}

// https://adlerhsieh.com/blog/write-your-own-go-linters-with-parser-package
// Visit function walks through each node in a file
func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	funcDecl, ok := n.(*ast.FuncDecl)
	if ok {
		fmt.Printf("Parsing func: [%v]...\n", funcDecl.Name)
	}

	blockStmt, ok := n.(*ast.BlockStmt)
	if ok {
		parseBlockStmt(blockStmt)
	}

	return v
}

func parseBlockStmt(stmt *ast.BlockStmt) bool {
	if len(stmt.List) < 4 {
		// do not need to parse
		return true
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

	if nps.isExchange() {
		fmt.Println(curfset.Position(stmt.Pos()))
	}

	return false
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
		fmt.Printf("WARN: Cannot convert right [%v]!\n",
			curfset.Position(assStmt.Pos()))
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

func (p *NamePairs) isExchange() bool {
	if len(p.arr) < 3 {
		return false
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
			return true
		}
	}

	return false
}
