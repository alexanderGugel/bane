[![Build Status](https://travis-ci.org/alexanderGugel/bane.svg?branch=master)](https://travis-ci.org/alexanderGugel/bane)

bane
====

Making the development of concurrent UDP based daemons less painful.

Example (Echo Server)
---------------------

```go
package main

import (
    "net"
    "github.com/alexanderGugel/bane"
)

func main() {
    addr, _ := net.ResolveUDPAddr("udp", "localhost:1337")
    conn, _ := net.ListenUDP("udp", addr)

    d := bane.New(conn, 1000)
    for {
        p := <-d.In
        d.Out <- p
    }
}
```