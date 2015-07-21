package bane // import "github.com/alexanderGugel/bane"

import (
	"net"
	"runtime"
)

// Packet represents an entity received a given client.
type Packet struct {
	Addr *net.UDPAddr // UDPAddr represents the address of the client's UDP end point.
	Data []byte       // Data represents the binary data that has been received by the client.
	Err  error        // Optional error associated with the packet, either encountered during receival or sending procedure.
}

// Daemon represents a UDP server listening that uses the specified UDP connection and exposes corresponding channels used for interacting with the network.
type Daemon struct {
	Conn       *net.UDPConn // The UDPConn connection that is being read from.
	Out        chan *Packet // Channel exposing outgoing packets.
	In         chan *Packet // Channel exposing incoming packets.
	OutErr     chan *Packet // Channel exposing errors encountered while writing to the UDP connection.
	InErr      chan *Packet // Channel exposing errors encountered while reading from the UDP connection.
	PacketSize int          // Maximum allowed packet size for incoming packets. Used for buffer initialisation.
}

func (d *Daemon) out() {
	for packet := range d.Out {
		_, err := d.Conn.WriteToUDP(packet.Data, packet.Addr)
		if err != nil {
			packet.Err = err
			d.OutErr <- packet
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
		} else {
			d.In <- &packet
		}
	}
}

// New returns a new bane daemon listening on the passed in UDP connection.
// packetSize specifies the maximum allowable packet size for incoming messages.
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

// NewFromAddr returns a new bane daemon bound to the specified network address.
// packetSize specifies the maximum allowable packet size for incoming messages.
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
