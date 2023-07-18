# DBuffer

[![Go Doc](https://pkg.go.dev/badge/github.com/ETZhangSX/dbuffer.svg)](https://pkg.go.dev/github.com/ETZhangSX/dbuffer)
[![Go Report Card](https://goreportcard.com/badge/github.com/ETZhangSX/dbuffer/v0.svg)](https://goreportcard.com/report/github.com/ETZhangSX/dbuffer)

Package dbuffer is a Go 1.18+ package using double-buffer to achieve data hot update.

Usage:

```go
package main

import (
    "fmt"
    "log"

    "github.com/ETZhangSX/dbuffer"
)

// implement loader for updating data.
type loader struct {
}

func(l *loader) Load(dest *map[string]string) {
    // update data
    *dest = make(map[string]string)
    *dest["key"] = "value"
}

func main() {
    // alloc func to create obj
    var alloc dbuffer.Alloc[map[string]string] = func() map[string]string {
        obj := map[string]string{}
        return obj
    }
    // update interval
    interval := 30 * time.Second
    buf := dbuffer.New[map[string]string](&loader{}, alloc, dbuffer.WithInterval(interval))

    // get data
    m1 := buf.Data()
    fmt.Println(m1["key"])

    // get data with done, which will block updates to current buffer util all refs done.
    m2, done := buf.DataWithDone()
    defer done()
    fmt.Println(m2["key"])
}
```
