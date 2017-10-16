# Gobstract
Simple automatic abstract text generator.

## Demo

```go
package main

import (
    "fmt"
    "time"
    "io/ioutil"
    "github.com/lucasmenendez/gobstract"
)

func main() {
    start := time.Now()

    var text string
    if raw, err := ioutil.ReadFile("demo/input"); err != nil {
        fmt.Println(err)
        return
    } else {
        text = string(raw)
    }

    if abstract, err := gobstract.NewAbstract(text, "es"); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("RESULTS\n")
        fmt.Println("Best sentence:")
        fmt.Println(abstract.GetBestSentence())

        fmt.Println("\nKeywords:")
        for _, s := range abstract.GetKeywords() {
            fmt.Printf("%s, ", s)
        }
        
        var length int = 0
        fmt.Println("\n\nSummary:")
        for _, sentence := range abstract.GetHightlights(10) {
            fmt.Println(sentence)
            length += len(sentence)
        }

        fmt.Println("\nReduced in:")
        fmt.Printf("Original: %d\n", len(text))
        fmt.Printf("Result: %d (%d%%)", length, (length * 100)/len(text))
    }

    elapsed := time.Since(start)
    fmt.Printf("\nTime elapsed: %s\n", elapsed)
}
```