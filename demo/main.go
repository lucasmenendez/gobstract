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

	var output float32
	for _, sentence := range abstract.Sentences {
		output += float32(len(sentence))
		fmt.Println(sentence)
	}

	var total float32 = float32(len(input))
	var percent float32 = output/total*100
	fmt.Printf("\nDeleted %g%% of content. %d sentences selected.\n", (100 - percent), len(abstract.Sentences))
}
