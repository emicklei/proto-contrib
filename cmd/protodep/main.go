// Copyright (c) 2018 Ernest Micklei
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
	"fmt"
	"log"
	"os"

	"github.com/emicklei/proto"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing proto file parameter")
	}
	root := os.Args[1]
	list, err := allImportsOf(root, map[string]bool{})
	if err != nil {
		log.Fatal("failed to parse imports ", err)
	}
	fmt.Println(root)
	for _, each := range list {
		fmt.Println(each.Filename)
	}
}

func allImportsOf(path string, visitMap map[string]bool) (list []*proto.Import, err error) {
	visitMap[path] = true
	imports, err := importsOf(path)
	if err != nil {
		return list, err
	}
	for _, each := range imports {
		if _, visited := visitMap[each.Filename]; !visited {
			list = append(list, each)
			sublist, err := allImportsOf(each.Filename, visitMap)
			if err != nil {
				return list, err
			}
			list = append(list, sublist...)
		}
	}
	return
}

func importsOf(path string) (list []*proto.Import, err error) {
	reader, err := os.Open(path)
	if err != nil {
		log.Println("failed to open", path)
		return list, err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return list, err
	}
	for _, each := range definition.Elements {
		if i, ok := each.(*proto.Import); ok {
			list = append(list, i)
		}
	}
	return
}
