package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseGoSrc(fname string) {
	bs, err := ReadBinFile(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", string(bs), parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	//printDecl(f, fset)

	// Print the AST.
	//ast.Print(fset, f) // https://zhuanlan.zhihu.com/p/28516587

	//递归调用逐一打印节点
	// https://www.jianshu.com/p/443bd82863f8
	ast.Inspect(f, func(n ast.Node) bool {
		ast.Print(fset, n)
		return true
	})

	//fmt.Printf("%#v\n", f)
	//var v visitor
	//ast.Walk(v, f)
}

// https://golang.google.cn/pkg/go/token/
func printDecl(f *ast.File, fset *token.FileSet) {
	// Print the location and kind of each declaration in f.
	for _, decl := range f.Decls {
		// Get the filename, line, and column back via the file set.
		// We get both the relative and absolute position.
		// The relative position is relative to the last line directive.
		// The absolute position is the exact position in the source.
		pos := decl.Pos()
		relPosition := fset.Position(pos)
		absPosition := fset.PositionFor(pos, false)

		// Either a FuncDecl or GenDecl, since we exit on error.
		kind := "func"
		if gen, ok := decl.(*ast.GenDecl); ok {
			kind = gen.Tok.String()
		}

		// If the relative and absolute positions differ, show both.
		fmtPosition := relPosition.String()
		if relPosition != absPosition {
			fmtPosition += "[" + absPosition.String() + "]"
		}

		fmt.Printf("%s: %s\n", fmtPosition, kind)
	}
}
