[![GoDoc](https://godoc.org/github.com/lucasmenendez/gobstract?status.svg)](https://godoc.org/github.com/lucasmenendez/gobstract)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasmenendez/gobstract)](https://goreportcard.com/report/github.com/lucasmenendez/gobstract)

# Gobstract
Gobstract package make extraction summaries from text provided. The algorithm measures sentence relations (measuring relevant token similarity), position and length to pick the text highlights.

## Installation
### PoS Tagging
For more information check instructions [here](https://github.com/lucasmenendez/gopostagger#train-corpus).

### Abstracts
```bash
export MODELS="<postagging trained models folder path>"

go get github.com/lucasmenendez/gobstact
```

### Use it
```go
package main

import (
    "fmt"
    "io/ioutil"
    "github.com/lucasmenendez/gobstract"
)

func main() {
    var input string
    if raw, err := ioutil.ReadFile("input"); err != nil {
        fmt.Println(err)
        return
    } else {
        input = string(raw)
    }

    if t, err := gobstract.NewText(input, "es"); err != nil {
        fmt.Println(err)
    } else {
        var summary []string = t.Summarize()
        for _, sentence := range summary {
            fmt.Println(sentence)
        }
    }    
}
```