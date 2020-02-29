package ping

import (
	"net"
	"time"

	libping "github.com/sparrc/go-ping"
)

type Pinger interface {
	GetLatency(address *net.UDPAddr) (latency time.Duration, err error)
}

type pinger struct{}

func NewPinger() Pinger {
	return &pinger{}
}

func (p *pinger) GetLatency(address *net.UDPAddr) (latency time.Duration, err error) {
	ipAddr := &net.IPAddr{
		IP:   address.IP,
		Zone: address.Zone,
	}
	pinger, err := libping.NewPinger(ipAddr.String())
	if err != nil {
		return 0, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.OnRecv = func(packet *libping.Packet) {
		latency = packet.Rtt
	}
	pinger.Run()
	return latency, nil
}
