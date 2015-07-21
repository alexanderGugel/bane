[![Build Status](https://travis-ci.org/alexanderGugel/bane.svg?branch=master)](https://travis-ci.org/alexanderGugel/bane)

bane
====

Making the development of concurrent UDP based daemons *fun*.

Example (Echo Server)
---------------------

```go
package main

import (
    "net"
    "github.com/alexanderGugel/bane"
)

func main() {
    d, _ := bane.NewFromAddr("udp", "localhost:1337", 1000)
    for {
        p := <-d.In
        d.Out <- p
    }
}
```
