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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/emicklei/proto"
)

var format = flag.String("format", "", "by default the output is a list of import names in plain text. JSON is the alternative")

func main() {
	flag.Parse()
	if len(flag.Args()) == 1 {
		log.Fatal("missing proto file parameter")
	}

	overallList := []string{}
	for _, each := range flag.Args() {
		overallList = append(overallList, each)
		list, err := allImportsOf(each, map[string]bool{})
		if err != nil {
			log.Fatal("failed to parse imports ", err)
		}
		for _, other := range list {
			overallList = append(overallList, other.Filename)
		}
	}

	overallList = unique(overallList)

	if len(*format) == 0 {
		for _, each := range overallList {
			fmt.Println(each)
		}
		return
	}
	if "json" == *format {
		json.NewEncoder(os.Stdout).Encode(overallList)
		return
	}
	log.Fatal("unknown format argument", *format)
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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
