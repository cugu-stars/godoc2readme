<h1 align="center">godoc2readme</h1>

<p  align="center">
 <a href="https://goreportcard.com/report/github.com/cugu/godoc2readme"><img src="https://goreportcard.com/badge/github.com/cugu/godoc2readme" alt="report" /></a>
 <a href="https://pkg.go.dev/github.com/cugu/godoc2readme"><img src="https://godoc.org/github.com/cugu/godoc2readme?status.svg" alt="doc" /></a>
</p>

The godoc2readme project can create README markdown files from Go package comments.

## Installation

```shell
go get -u github.com/cugu/godoc2readme
```

## Usage
Just point godoc2readme to the directory with the go.mod file.

```
godoc2readme my/gopackage > README.md
```

godoc2readme can also process custom templates. Those custom templates follow the Go template
syntax ([https://golang.org/pkg/text/template/](https://golang.org/pkg/text/template/)). All fields available in the GoPackage struct are
available in the template.

```
godoc2readme --template mycusomreadme.tmpl.md . > README.md
```
## Example Template
The following shows an example template that prints a package and all its subpackages.

```
# {{.Name}} ({{.ModulePath}})

{{.Doc}}

{{if .Subpackages}}{{ range $key, $value := .Subpackages }}{{if $value.Doc }}
## {{ $value.Name }} {{if .Command}}command{{else}}library{{end}}
{{ $value.Doc }}
{{end}}{{end}}{{end}}
```







