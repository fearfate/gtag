package gtag

import (
	"bytes"
	"text/template"
)

type templateData struct {
	Package string
	Command string
	Types   []templateDataType
	Tags    []templateDataTag
}

type templateDataType struct {
	Name   string
	Fields []templateDataTypeField
}

type templateDataTypeField struct {
	Name string
	Tag  string
}

type templateDataTag struct {
	Name  string
	Value string
}

const templateLayout = `
// Code generated by gtag. DO NOT EDIT.
// See: https://github.com/wolfogre/gtag

//go:generate {{.Command}}
package {{.Package}}

import (
	"reflect"
	"strings"
)

{{$tags := .Tags}}
{{- range .Types}}

var (
	valueOf{{.Name}} = {{.Name}}{}
	typeOf{{.Name}}  = reflect.TypeOf(valueOf{{.Name}})

{{$type := .Name}}
{{- range .Fields}}
	_ = valueOf{{$type}}.{{.Name}}
	fieldOf{{$type}}{{.Name}}, _ = typeOf{{$type}}.FieldByName("{{.Name}}")
	tagOf{{$type}}{{.Name}} = fieldOf{{$type}}{{.Name}}.Tag
{{end}}
)

// {{$type}}Tags indicate tags of type {{$type}}
type {{$type}}Tags struct {
{{- range .Fields}}
	{{.Name}} string // {{.Tag}}
{{- end}}
}

// Tags return specified tags of {{$type}}
func (*{{$type}}) Tags(tag string, convert ...func(string) string) {{$type}}Tags {
	conv := func(in string) string { return strings.TrimSpace(strings.Split(in, ",")[0]) }
	if len(convert) > 0 {
		conv = convert[0]
	}
	if conv == nil {
		conv = func(in string) string { return in }
	}
	return {{$type}}Tags{
{{- range .Fields}}
		{{.Name}}: conv(tagOf{{$type}}{{.Name}}.Get(tag)),
{{- end}}
	}
}

// TagSlice return specified tag slice of {{$type}}
func (*{{$type}}) TagSlice(tag string, convert ...func(string) string) []string {
	conv := func(in string) string { return strings.TrimSpace(strings.Split(in, ",")[0]) }
	if len(convert) > 0 {
		conv = convert[0]
	}
	if conv == nil {
		conv = func(in string) string { return in }
	}

	return []string{
{{- range .Fields}}
		conv(tagOf{{$type}}{{.Name}}.Get(tag)),
{{- end}}
	}
}

{{range $tags}}
// Tags{{.Name}} is alias of Tags("{{.Value}}")
func (*{{$type}}) Tags{{.Name}}() {{$type}}Tags {
	var v *{{$type}}
	return v.Tags("{{.Value}}")
}

// TagSlice{{.Name}} is alias of TagSlice("{{.Value}}")
func (*{{$type}}) TagSlice{{.Name}}() []string {
	var v *{{$type}}
	return v.TagSlice("{{.Value}}")
}

{{end}}

{{- end}}

`

func execute(data templateData) []byte {
	tp := template.Must(template.New("").Parse(templateLayout))

	out := &bytes.Buffer{}
	if err := tp.Execute(out, data); err != nil {
		panic(err)
	}

	return out.Bytes()
}
