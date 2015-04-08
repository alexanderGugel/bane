package bane

import (
	"net"
	"runtime"
)

type Packet struct {
	Addr *net.UDPAddr
	Data []byte
	Err  error
}

type Daemon struct {
	Conn       *net.UDPConn
	Out        chan *Packet
	In         chan *Packet
	OutErr     chan *Packet
	InErr      chan *Packet
	PacketSize int
}

func (d *Daemon) out() {
	for packet := range d.Out {
		_, err := d.Conn.WriteToUDP(packet.Data, packet.Addr)
		if err != nil {
			packet.Err = err
			d.OutErr <- packet
			continue
		}
	}
}

func (d *Daemon) in() {
	for {
		b := make([]byte, d.PacketSize)
		n, addr, err := d.Conn.ReadFromUDP(b)
		packet := Packet{addr, b[:n], err}
		if err != nil {
			d.InErr <- &packet
			continue
		}
		d.In <- &packet
	}
}

func New(conn *net.UDPConn, packetSize int) *Daemon {
	d := &Daemon{
		Conn:       conn,
		In:         make(chan *Packet),
		Out:        make(chan *Packet),
		OutErr:     make(chan *Packet),
		InErr:      make(chan *Packet),
		PacketSize: packetSize,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go d.in()
		go d.out()
	}

	return d
}

func NewFromAddr(network string, addr string, packetSize int) (*Daemon, error) {
	resolvedAddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP(network, resolvedAddr)
	if err != nil {
		return nil, err
	}
	return New(conn, packetSize), nil
}
