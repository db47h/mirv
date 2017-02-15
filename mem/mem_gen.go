// Copyright © 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build gen

package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

var license = `{{define "license"}}// Copyright © 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.{{end}}`

var tpl = `
{{- template "license"}}

package mem

import (
	"encoding/binary"

	"github.com/db47h/mirv"
)
{{define "T1" -}}
	len((*m)[addr:]) < {{.Bytes}}
{{- end -}}

{{range .}}{{$od := .OD}}{{$of := .OF}}{{$on := .ON}}
// {{$od}} memory interface
type {{$on}} []uint8

func (m *{{$on}}) Size() mirv.Address { return mirv.Address(len(*m)) }

func (m *{{$on}}) Type() Type { return MemRAM }

func (m *{{$on}}) ByteOrder() mirv.ByteOrder { return mirv.{{.OF}} }
{{range .Width}}
// Read{{.Bits}} returns the {{.Bits}} bits {{$od}} value at address addr.
//
func (m *{{$on}}) Read{{.Bits}}(addr mirv.Address) (uint{{.Bits}}, error) {
	if {{template "T1" .}} {
		return 0, errPage
	}
	{{if (eq .Bits 8) -}}
	return (*m)[addr], nil
	{{- else -}}
	return binary.{{$of}}.Uint{{.Bits}}((*m)[addr:]), nil
	{{- end}}
}

// Write{{.Bits}} writes the {{.Bits}} bits {{$od}} value to address addr.
//
func (m *{{$on}}) Write{{.Bits}}(addr mirv.Address, v uint{{.Bits}}) error {
	if {{template "T1" .}} {
		return errPage
	}
	{{if (eq .Bits 8) -}}
	(*m)[addr] = v
	{{- else -}}
	binary.{{$of}}.PutUint{{.Bits}}((*m)[addr:], v)
	{{- end}}
	return nil
}
{{end}}{{end}}
{{- define "noMemory" -}}
{{range .}}
// Read{{.Bits}} always returns 0 and an error of type *ErrBus.
//
func (NoMemory) Read{{.Bits}}(addr mirv.Address) (uint{{.Bits}}, error) {
	return 0, errBus(opRead, {{.Bytes}}, addr)
}

// Write{{.Bits}} always returns an error of type *ErrBus.
//
func (NoMemory) Write{{.Bits}}(addr mirv.Address, v uint{{.Bits}}) error {
	return errBus(opWrite, {{.Bytes}}, addr)
}
{{end}}{{end}}`

type width struct {
	Bits  int
	Bytes int
}

type data struct {
	OD    string
	ON    string
	OF    string
	Width []width
}

var out = flag.String("o", "", "output file")

func main() {
	flag.Parse()

	if *out == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t := template.New("tpl")
	t, _ = t.Parse(license)
	t, err = t.Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	if err = t.Execute(f, []data{
		{
			OD:    "big endian ",
			ON:    "bigEndian",
			OF:    "BigEndian",
			Width: []width{{8, 1}, {16, 2}, {32, 4}, {64, 8}},
		},
		{
			OD:    "little endian ",
			ON:    "littleEndian",
			OF:    "LittleEndian",
			Width: []width{{8, 1}, {16, 2}, {32, 4}, {64, 8}},
		},
	}); err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(f, "noMemory", []width{{8, 1}, {16, 2}, {32, 4}, {64, 8}})
	// t = template.New("tpl2")
	// t, err = t.Parse(noMem)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// t.
}
