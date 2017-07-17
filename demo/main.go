package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lucasmenendez/gobstract"
)

var language = "es"
var input_path = "./demo/input_es"

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
	if abstract, err = gobstract.NewAbstract(input, language); err != nil {
		panic(err)
	}

	var output int
	for _, sentence := range abstract.Sentences {
		output += len(sentence)
		fmt.Println(sentence)
	}

	var total float32 = float32(len(input))
	var percent float32 = float32(output) / total * 100.0
	fmt.Printf("\n\tDeleted:\t%d/%d (%g%%).\n", len(input)-output, len(input), (100.0 - percent))
	fmt.Printf("\tSelected:\t%d sentences.\n\n", len(abstract.Sentences))
}
