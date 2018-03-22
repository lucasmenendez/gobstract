# Gobstract
Gobstract package make extraction summaries from text provided. The algorithm measures sentence relations (measuring relevant token similarity), position and length to pick the text highlights.
## Demo

```go
package main

import (
    "fmt"
    "io/ioutil"
    "github.com/lucasmenendez/gobstract"
)

func main() {
    var input string
    if raw, err := ioutil.ReadFile("demo/input"); err != nil {
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