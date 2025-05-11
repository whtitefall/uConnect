package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Chunk struct {
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Code      string `json:"code"`
}

func main() {
	root := "./forum_backend"
	fset := token.NewFileSet()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
			return nil
		}

		if filepath.Ext(path) == ".go" {
			src, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read error: %v\n", err)
				return nil
			}

			node, err := parser.ParseFile(fset, path, src, parser.ParseComments)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
				return nil
			}

			for _, decl := range node.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					start := fset.Position(fn.Pos()).Offset
					end := fset.Position(fn.End()).Offset

					if start >= 0 && end <= len(src) && start < end {
						chunk := Chunk{
							File:      path,
							Symbol:    fn.Name.Name,
							StartLine: fset.Position(fn.Pos()).Line,
							EndLine:   fset.Position(fn.End()).Line,
							Code:      string(src[start:end]),
						}

						jsonChunk, _ := json.Marshal(chunk)
						fmt.Println(string(jsonChunk))
					}
				}
			}
		}
		return nil
	})
}
