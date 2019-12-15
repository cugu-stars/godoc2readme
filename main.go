// Copyright (c) 2019 Jonas Plum
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

// Package godoc2readme can create README markdown files from Go package comments.
//
// Usage
//
// Just point godoc2readme to the directory with the go.mod file.
//
//     godoc2readme my/gopackage > README.md
//
//
// godoc2readme can also process custom templates. Those custom templates follow the Go template
// syntax (https://golang.org/pkg/text/template/). All fields available in the GoPackage struct are
// available in the template.
//
//     godoc2readme --template mycusomreadme.tmpl.md . > README.md
//
// Example Template
//
// The following shows an example template that prints a package and all its subpackages.
//
//     # {{.Name}} ({{.ModulePath}})
//
//     {{.Doc}}
//
//     {{if .Subpackages}}{{ range $key, $value := .Subpackages }}{{if $value.Doc }}
//     ## {{ $value.Name }} {{if .Command}}command{{else}}library{{end}}
//     {{ $value.Doc }}
//     {{end}}{{end}}{{end}}
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/doc"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/cugu/godoc2readme/templates"
)

// GoPackage contains fields that can be used in templates to create markdown files for go packages.
type GoPackage struct {
	Name          string
	ModulePath    string
	RelModulePath string

	Doc         string
	MainDoc     string
	Synopsis    string
	Command     bool
	Bugs        []string
	Subpackages map[string]*GoPackage
	Examples    map[string]string
	Code        string
}

func packageMarkdownDoc(diskPath string, goModPath string, recurse bool) *GoPackage {
	// get doc string
	bpkg, err := build.ImportDir(diskPath, build.ImportComment)
	logError(err)
	fset, docPkg, codeExamples, err := computeDoc(bpkg)
	logError(err)

	// generate examples
	examples := map[string]string{}
	if recurse {
		for _, example := range codeExamples {
			if example.Play == nil {
				continue
			}
			var code bytes.Buffer
			name := example.Name
			name = strings.ReplaceAll(name, "_", ".")
			format.Node(&code, fset, example.Play)

			examples[name] = code.String()
		}
	}

	// convert to markdown
	buf := &bytes.Buffer{}
	re := regexp.MustCompile(`^Package (\w+)`)
	niceDoc := re.ReplaceAllString(docPkg.Doc, "The $1 project")
	toMarkdown(buf, niceDoc, nil)

	sBuf := &bytes.Buffer{}
	syn := doc.Synopsis(niceDoc)
	niceSDoc := niceDoc[len(syn):]
	toMarkdown(sBuf, niceSDoc, nil)

	_, base := path.Split(goModPath)

	// get subprojects
	subPkg := map[string]*GoPackage{}
	if recurse {
		err = filepath.Walk(diskPath, func(subpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && (info.Name() == "assets" || info.Name() == ".git" || info.Name() == "node_modules") {
				return filepath.SkipDir
			}
			if info.IsDir() {
				subbpkg, err := build.ImportDir(subpath, build.ImportComment)
				if err == nil {
					if subbpkg.Name == base {
						return nil
					}

					subMod, err := filepath.Rel(diskPath, subpath)
					if err != nil {
						return err
					}
					subMod = filepath.ToSlash(subMod)
					subPackage := packageMarkdownDoc(subpath, path.Join(goModPath, subMod), false)
					if subPackage.Name != base && subPackage.Doc != "" {
						subPkg[subPackage.Name] = subPackage
					}
				}
			}
			return nil
		})
		logError(err)
	}
	if len(subPkg) == 0 {
		subPkg = nil
	}

	return &GoPackage{
		Name:          base,
		Doc:           buf.String(),
		MainDoc:       sBuf.String(),
		ModulePath:    goModPath,
		Synopsis:      syn,
		RelModulePath: strings.ReplaceAll(goModPath, "github.com/", ""),
		Command:       docPkg.Name == "main",
		Bugs:          docPkg.Bugs,
		Subpackages:   subPkg,
		Examples:      examples,
	}
}

func main() {
	templatePathP := flag.String("template", "Readme", "path of template file")
	flag.Parse()
	templatePath := *templatePathP

	for _, modPath := range flag.Args() {
		// get module path
		b, err := ioutil.ReadFile(filepath.Join(modPath, "go.mod"))
		logError(err)
		modulePath := modulePath(b)

		p := packageMarkdownDoc(modPath, modulePath, true)

		// parse template
		var tmpl string
		_, err = os.Stat(templatePath)
		if err != nil {
			if value, ok := templates.Templates[templatePath]; ok && os.IsNotExist(err) {
				tmpl = value
			} else {
				logError(fmt.Errorf("Template %s does not exist", templatePath))
			}
		} else {
			tmplBytes, err := ioutil.ReadFile(templatePath)
			logError(err)
			tmpl = string(tmplBytes)
		}

		t, err := template.New("doc").Parse(tmpl)
		logError(err)

		// print template
		err = t.ExecuteTemplate(os.Stdout, "doc", p)
		logError(err)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
