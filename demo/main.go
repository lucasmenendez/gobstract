package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lucasmenendez/gobstract"
)

const language string = "es"

func main() {
	var err error
	var inputPath string = fmt.Sprintf("./demo/input_%s", language)

	var inputRoot string
	if inputRoot, err = filepath.Abs(inputPath); err != nil {
		panic(err)
	}

	var inputRaw []byte
	if inputRaw, err = ioutil.ReadFile(inputRoot); err != nil {
		panic(err)
	}

	var input string = string(inputRaw)
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
