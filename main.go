package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-test-renamer <directory>")
		os.Exit(1)
	}

	directory := os.Args[1]

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		if err := processFile(path); err != nil {
			fmt.Printf("Error processing file %s: %s\n", path, err)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error processing files: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("All files processed successfully.")
}

func processFile(filename string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %s", err)
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

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err)
	}
	defer out.Close()

	if err := format.Node(out, fset, node); err != nil {
		return fmt.Errorf("error writing to output file: %s", err)
	}

	fmt.Printf("File processed: %s\n", filename)

	return nil
}
