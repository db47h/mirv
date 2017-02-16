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
	"github.com/db47h/mirv"
)
{{define "T1" -}}
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
{{- end -}}
{{range .}}
// Read{{.}} returns the {{.}} bits value at address addr.
//
func (b *Bus) Read{{.}}(addr mirv.Address) (uint{{.}}, error) {
	{{template "T1"}}
	return blk.m.Read{{.}}(addr - blk.s)
}

// Write{{.}} writes the {{.}} bits value to address addr.
//
func (b *Bus) Write{{.}}(addr mirv.Address, v uint{{.}}) error {
	{{template "T1"}}
	return blk.m.Write{{.}}(addr-blk.s, v)
}
{{end}}`

type data struct {
	OD    string
	ON    string
	Width []int
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
	if err := t.Execute(f, []int{8, 16, 32, 64}); err != nil {
		log.Fatal(err)
	}
}
