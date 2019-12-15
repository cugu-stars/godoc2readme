// Copyright 2019 Dmitri Shuralyov
// Copyright 2019 Jonas Plum
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// This code was adapted from the original
// https://github.com/shurcooL/cmd/blob/master/gorepogen/main.go.

package main

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
)

// computeDoc computes the package documentation for the given package.
func computeDoc(bpkg *build.Package) (*token.FileSet, *doc.Package, []*doc.Example, error) {
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	var filelist []*ast.File

	var gofiles []string
	gofiles = append(gofiles, bpkg.CgoFiles...)
	gofiles = append(gofiles, bpkg.GoFiles...)
	gofiles = append(gofiles, bpkg.TestGoFiles...)
	gofiles = append(gofiles, bpkg.XTestGoFiles...)

	for _, file := range gofiles {
		f, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, file), nil, parser.ParseComments)
		if err != nil {
			return nil, nil, nil, err
		}
		files[file] = f
		filelist = append(filelist, f)
	}

	examples := doc.Examples(filelist...)

	apkg := &ast.Package{
		Name:  bpkg.Name,
		Files: files,
	}
	return fset, doc.New(apkg, bpkg.ImportPath, 0), examples, nil
}
