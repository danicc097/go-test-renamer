package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {
	var verbose bool
	excludedDirs := []string{"vendor"}

	ed := flag.String("exclude", "", "Comma-separated directories to exclude")
	v := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	if ed != nil {
		excludedDirs = append(excludedDirs, strings.Split(*ed, ",")...)
	}

	if v != nil {
		verbose = *v
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && slices.Contains(excludedDirs, info.Name()) {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, "_test.go") && !info.IsDir() {
			input, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("Error reading file %s: %s\n", path, err)
			}

			output, err := processFile(bytes.NewReader(input))
			if err != nil {
				log.Fatalf("Error processing file %s: %s\n", path, err)
			}

			err = os.WriteFile(path, output, info.Mode())
			if err != nil {
				log.Fatalf("Error writing processed contents to file %s: %s\n", path, err)
			}

			if verbose {
				fmt.Printf("Processed file %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error during file processing: %s\n", err)
	}
}

func processFile(input io.Reader) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", input, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	_, err = conf.Check("", fset, []*ast.File{node}, info)
	if err != nil {
		return nil, fmt.Errorf("type checking error: %s", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if xIdent, ok := selExpr.X.(*ast.Ident); ok {
					if xIdent.Name == "t" && selExpr.Sel.Name == "Run" {
						tRunArgs := callExpr.Args
						if len(tRunArgs) > 0 {
							testNameArg := tRunArgs[0]

							// update directly
							if lit, ok := testNameArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
								lit.Value = strings.ReplaceAll(lit.Value, " ", "_")
							}

							// handle test case name via table driven test
							if sel, ok := testNameArg.(*ast.SelectorExpr); ok {
								if sel.Sel.Name == "name" || sel.Sel.Name == "Name" {
									if obj := info.ObjectOf(sel.Sel); obj != nil {
										fmt.Printf("obj: %v\n", obj)
										fmt.Printf("sel.Sel: %+v\n", info.ObjectOf(sel.Sel))

										if tv, ok := obj.(*types.Var); ok && tv.IsField() {
											fmt.Printf("sel.X: %T\n", sel.X)

											// TODO: update referenced

											// ast.Inspect(node, func(n ast.Node) bool {
											// 	if ident, ok := n.(*ast.Ident); ok {
											// 		if ident.Name == obj.Name() {
											// 			if ident.Obj != nil {
											// 				fmt.Printf("ident.Decl: %+v\n", ident.Obj.Decl)
											// 			}
											// 			// ...
											// 		}
											// 	}
											// 	return true
											// })
										}
									}
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
