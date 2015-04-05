package bane

import (
    "testing"
    "net"
)

func TestEcho(t *testing.T) {
    addr, err := net.ResolveUDPAddr("udp", "localhost:0")
    if err != nil {
        t.Fatal(err)
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        t.Fatal(err)
    }

    d := New(conn, 1000)

    realAddr := conn.LocalAddr().(*net.UDPAddr)

    d.Out <- &Packet{realAddr, []byte("TestEcho"), nil}

    p := <- d.In
    if string(p.Data) != "TestEcho" {
        t.Fatalf("want %v, got %v", string(p.Data), "TestEcho")
    }
    p.Data = []byte("TestEchoResponse")
    d.Out <- p
    if string(p.Data) != "TestEchoResponse" {
        t.Fatalf("want %v, got %v", string(p.Data), "TestEchoResponse")
    }
}

func BenchmarkEcho(b *testing.B) {
    addr, err := net.ResolveUDPAddr("udp", "localhost:0")
    if err != nil {
        b.Fatal(err)
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        b.Fatal(err)
    }

    d := New(conn, 1000)

    realAddr := conn.LocalAddr().(*net.UDPAddr)


    go func() {
        for {
            d.Out <- &Packet{realAddr, []byte("BenchmarkEcho"), nil}
        }
    }()

    n := 0
    for n < b.N {
        n++
        p := <- d.In
        if string(p.Data) != "BenchmarkEcho" {
            b.Fatalf("want %v, got %v", string(p.Data), "BenchmarkEcho")
        }
        d.Out <- p
    }
}
