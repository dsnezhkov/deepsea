package main

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"log"
	"os"


	"jaytaylor.com/html2text"
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
	text, err := html2text.FromReader(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dstPath, []byte(text), 0644)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(text)

}

