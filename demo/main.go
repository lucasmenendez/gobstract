package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lucasmenendez/gobstract"
)

var input_path = "./demo/input"

func main() {
	var err error

	var input_root string
	if input_root, err = filepath.Abs(input_path); err != nil {
		panic(err)
	}

	var input_raw []byte
	if input_raw, err = ioutil.ReadFile(input_root); err != nil {
		panic(err)
	}

	var input string = string(input_raw)
	var abstract *gobstract.Gobstract
	if abstract, err = gobstract.NewAbstract(input, "en"); err != nil {
		panic(err)
	}


	for _, sentence := range abstract.Sentences {
		fmt.Println(sentence)
	}
}
