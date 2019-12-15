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

package templates

var readme = `<h1 align="center">{{.Name}}</h1>

<p  align="center">
 <a href="https://{{.ModulePath}}/actions"><img src="https://{{.ModulePath}}/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/{{.RelModulePath}}"><img src="https://codecov.io/gh/{{.RelModulePath}}/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/{{.ModulePath}}"><img src="https://goreportcard.com/badge/{{.ModulePath}}" alt="report" /></a>
 <a href="https://pkg.go.dev/{{.ModulePath}}"><img src="https://godoc.org/{{.ModulePath}}?status.svg" alt="doc" /></a>
</p>

{{.Synopsis}}

## Installation

` + "```" + `shell
go get -u {{.ModulePath}}
` + "```" + `

{{.MainDoc}}

{{if .Examples}}
## Examples
{{ range $key, $value := .Examples }}

{{if $key}}### {{ $key }}{{end}}
` + "```" + ` go
{{ $value }}
` + "```" + `
{{end}}{{end}}

{{if .Bugs}}
## Known Bugs
{{range .Bugs}}* {{.}}{{end}}
{{end}}

{{if .Subpackages}}
## Subpackages

| Package | Description |
| --- | --- |
{{ range $key, $value := .Subpackages }}| **{{ $value.Name }}** | {{ $value.Synopsis }} |
{{end}}{{end}}
`
