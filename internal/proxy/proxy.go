package proxy

import (
	"context"
	"net"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/pingodown/internal/connection"
	"github.com/qdm12/pingodown/internal/ping"
	"github.com/qdm12/pingodown/internal/state"
)

type Proxy interface {
	Run(ctx context.Context) error
}

type proxy struct {
	bufferSize    int
	proxyConn     *net.UDPConn
	serverAddress *net.UDPAddr
	state         state.State
	logger        logging.Logger
	pinger        ping.Pinger
}

func NewProxy(listenAddress, serverAddress string, logger logging.Logger) (Proxy, error) {
	state := state.NewState()
	p := &proxy{
		bufferSize: 65535,
		state:      state,
		logger:     logger,
		pinger:     ping.NewPinger(),
	}
	var err error
	proxyAddress, err := net.ResolveUDPAddr("udp", listenAddress)
	if err != nil {
		return nil, err
	}
	p.proxyConn, err = net.ListenUDP("udp", proxyAddress)
	if err != nil {
		return nil, err
	}
	p.serverAddress, err = net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *proxy) Close() error {
	return p.proxyConn.Close()
}

type clientPacket struct {
	clientAddress *net.UDPAddr
	data          []byte
}

func pingClients(pinger ping.Pinger, state state.State, logger logging.Logger) {
	for _, address := range state.GetClientAddresses() {
		address := address
		latency, err := pinger.GetLatency(&address)
		if err != nil {
			logger.Error(err)
			continue
		}
		state.SetLatency(&address, latency)
	}
}

func updatePings(state state.State, logger logging.Logger) {
	highestLatency := state.GetHighestLatency()
	for _, address := range state.GetClientAddresses() {
		address := address
		conn, err := state.GetConnection(&address)
		if err != nil {
			logger.Error(err)
			continue
		}
		currentLatency, err := state.GetLatency(&address)
		if err != nil {
			logger.Error(err)
			continue
		}
		ping := highestLatency - currentLatency
		if !conn.PingApproximatesTo(ping) {
			logger.Info("Setting new ping of %s to %dms", conn.GetClientAddress(), ping.Milliseconds())
			conn.SetPing(ping) // works by pointer
		}
	}
}

func updatePingPeriodically(ctx context.Context, period time.Duration,
	pinger ping.Pinger, state state.State, logger logging.Logger) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			pingClients(pinger, state, logger)
			updatePings(state, logger)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (p *proxy) Run(ctx context.Context) (err error) {
	p.logger.Info("Running proxy to %s on %s", p.serverAddress, p.proxyConn.LocalAddr())
	packets := make(chan clientPacket, 100)
	go updatePingPeriodically(ctx, 10*time.Second, p.pinger, p.state, p.logger)
	go func() {
		if err := readFromClients(p.proxyConn, packets, p.bufferSize); err != nil {
			p.logger.Error(err)
		}
	}()
	for {
		select {
		case packet := <-packets:
			conn, err := p.state.GetConnection(packet.clientAddress)
			if err != nil {
				p.logger.Info("New client %s connecting", packet.clientAddress)
				conn, err = connection.NewConnection(p.serverAddress, packet.clientAddress, p.bufferSize)
				if err != nil {
					p.logger.Error(err)
					continue
				}
				conn = p.state.SetConnection(conn)
				go conn.ForwardServerToClient(ctx, p.proxyConn, p.logger)
			}
			go func() {
				if err := conn.WriteToServerWithDelay(ctx, packet.data); err != nil {
					p.logger.Error(err)
				}
			}()
		case <-ctx.Done():
			p.logger.Info("context canceled, closing proxy connection")
			return p.Close()
		}
	}
}

func readFromClients(proxy *net.UDPConn, packets chan<- clientPacket, bufferSize int) error {
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, clientAddress, err := proxy.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		data := make([]byte, bytesRead)
		copy(data, buffer[:bytesRead])
		packets <- clientPacket{
			clientAddress: clientAddress,
			data:          data,
		}
	}
}
