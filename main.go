package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		return
	}

	output, err := processFile(bytes.NewReader(input))
	if err != nil {
		fmt.Printf("Error processing file: %s\n", err)
		return
	}

	_, err = os.Stdout.Write(output)
	if err != nil {
		fmt.Printf("Error writing output: %s\n", err)
	}
}

func processFile(input io.Reader) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", input, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if xIdent, ok := selExpr.X.(*ast.Ident); ok {
					if xIdent.Name == "t" && selExpr.Sel.Name == "Run" {
						args := callExpr.Args
						if len(args) > 0 {
							if lit, ok := args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								oldName := lit.Value
								newName := strings.ReplaceAll(oldName, " ", "_")
								if newName != oldName {
									lit.Value = newName
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	var output bytes.Buffer
	err = format.Node(&output, fset, node)
	if err != nil {
		return nil, fmt.Errorf("error formatting output: %s", err)
	}

	return output.Bytes(), nil
}
