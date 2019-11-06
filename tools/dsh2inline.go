package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"os"


	"github.com/aymerick/douceur/inliner"
)

func usage(name string){
	fmt.Printf("Usage: %s <srcfile> <dstfile>\n", name)
	os.Exit(1)
}

func main() {

	if len(os.Args[1:]) != 2  {
		usage(os.Args[0])
	}
	var srcPath = os.Args[1]
	var dstPath = os.Args[2]

	bs, err := ioutil.ReadFile(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	html, err := inliner.Inline(string(bs))
	if err != nil {
		panic("Unable to inline")
	}

	// fmt.Println(html)
	err = ioutil.WriteFile(dstPath, []byte(html), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

