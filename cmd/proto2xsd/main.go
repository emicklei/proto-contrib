// Copyright (c) 2017 Ernest Micklei
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"io"
	"os"

	"flag"

	"bytes"

	xsd "github.com/emicklei/proto-contrib/pkg/proto2xsd"
)

var (
	oNamespace = flag.String("ns", "http://your.company.com/domain/version", "namespace of the target types")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	exitCode := 0
	for _, each := range flag.Args() {
		if err := readConvertWrite(each); err != nil {
			println(err.Error())
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func readConvertWrite(filename string) error {
	// open for read
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// buffer before write
	buf := new(bytes.Buffer)
	if err := xsd.Convert(filename, *oNamespace, file, buf); err != nil {
		return err
	}
	// write to stdout
	if _, err := io.Copy(os.Stdout, bytes.NewReader(buf.Bytes())); err != nil {
		return err
	}
	return nil
}
