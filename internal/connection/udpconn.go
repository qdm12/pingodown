package connection

import (
	"net"
)

//go:generate mockgen -destination=../../mocks/udpconn.go -package=mocks github.com/qdm12/pingodown/internal/connection UDPConn

type UDPConn interface {
	Close() error
	Read(b []byte) (read int, err error)
	Write(b []byte) (written int, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (written int, err error)
}

type udpConn struct {
	conn *net.UDPConn
}

func NewUDPConn(address *net.UDPAddr) (conn UDPConn, err error) {
	c, err := net.DialUDP("udp", nil, address)
	if err != nil {
		return nil, err
	}
	return &udpConn{
		conn: c,
	}, nil
}

func (u *udpConn) Close() error {
	return u.conn.Close()
}

func (u *udpConn) Read(b []byte) (read int, err error) {
	return u.conn.Read(b)
}

func (u *udpConn) Write(b []byte) (written int, err error) {
	return u.conn.Write(b)
}
func (u *udpConn) WriteToUDP(b []byte, addr *net.UDPAddr) (written int, err error) {
	return u.conn.WriteToUDP(b, addr)
}
