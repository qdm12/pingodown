package state

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/qdm12/pingodown/internal/connection"
)

type State interface {
	GetClientAddresses() (clientAddresses []*net.UDPAddr)
	// SetConnection sets a connection in the state and returns the saved connection or the
	// already existing connection as it does not overwrite an existing connection
	SetConnection(conn connection.Connection) connection.Connection
	// GetConnection retrieves an existing connection from the state
	GetConnection(clientAddress *net.UDPAddr) (conn connection.Connection, err error)
	SetLatency(clientAddress *net.UDPAddr, latency time.Duration)
	GetLatency(clientAddress *net.UDPAddr) (latency time.Duration, err error)
	GetHighestLatency() time.Duration
}

type state struct {
	// Key is the client IP address
	latencies      map[string]time.Duration
	latenciesMutex sync.RWMutex
	// Key is the client address
	connections      map[string]connection.Connection
	connectionsMutex sync.RWMutex
}

func NewState() State {
	return &state{
		connections: make(map[string]connection.Connection),
		latencies:   make(map[string]time.Duration),
	}
}

func (s *state) GetClientAddresses() (clientAddresses []*net.UDPAddr) {
	s.connectionsMutex.RLock()
	defer s.connectionsMutex.RUnlock()
	for _, conn := range s.connections {
		clientAddresses = append(clientAddresses, conn.GetClientUDPAddress())
	}
	return clientAddresses
}

func (s *state) GetConnection(clientAddress *net.UDPAddr) (conn connection.Connection, err error) {
	s.connectionsMutex.RLock()
	defer s.connectionsMutex.RUnlock()
	key := clientAddress.String()
	conn, ok := s.connections[key]
	if !ok {
		return nil, fmt.Errorf("no connection found for client address %s", key)
	}
	return conn, nil
}

func (s *state) SetConnection(conn connection.Connection) connection.Connection {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()
	key := conn.GetClientUDPAddress().String()
	if conn, ok := s.connections[key]; ok { // in case it still got created
		return conn
	}
	s.connections[key] = conn
	return conn
}

func (s *state) GetLatency(clientAddress *net.UDPAddr) (latency time.Duration, err error) {
	s.latenciesMutex.RLock()
	defer s.latenciesMutex.RUnlock()
	key := clientAddress.String()
	latency, ok := s.latencies[key]
	if !ok {
		return 0, fmt.Errorf("no latency found for client address %s", key)
	}
	return latency, nil
}

func (s *state) GetHighestLatency() (maxLatency time.Duration) {
	s.latenciesMutex.RLock()
	defer s.latenciesMutex.RUnlock()
	for _, latency := range s.latencies {
		if latency > maxLatency {
			maxLatency = latency
		}
	}
	return maxLatency
}

func (s *state) SetLatency(clientAddress *net.UDPAddr, latency time.Duration) {
	s.latenciesMutex.Lock()
	defer s.latenciesMutex.Unlock()
	key := clientAddress.String()
	s.latencies[key] = latency
}
