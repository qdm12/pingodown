package connection

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Connection interface {
	GetClientUDPAddress() *net.UDPAddr
	ForwardServerToClient(ctx context.Context, proxy *net.UDPConn, logger logging.Logger)
	WriteToServerWithDelay(ctx context.Context, data []byte) error
	SetPing(ping time.Duration)
}

// Information maintained for each client/server connection
type connection struct {
	clientAddress *net.UDPAddr
	server        UDPConn
	bufferSize    int
	inboundDelay  time.Duration
	outboundDelay time.Duration
	sync.RWMutex
}

// Generate a new connection by opening a UDP connection to the server
func NewConnection(serverAddress, clientAddress *net.UDPAddr, bufferSize int) (Connection, error) {
	server, err := NewUDPConn(serverAddress)
	if err != nil {
		return nil, err
	}
	return &connection{
		clientAddress: clientAddress,
		server:        server,
		bufferSize:    bufferSize,
	}, nil
}

func (c *connection) close() error {
	return c.server.Close()
}

func (c *connection) getInboundDelay() time.Duration {
	c.RLock()
	defer c.RUnlock()
	return c.inboundDelay
}

func (c *connection) getOutboundDelay() time.Duration {
	c.RLock()
	defer c.RUnlock()
	return c.outboundDelay
}

func (c *connection) SetPing(ping time.Duration) {
	c.Lock()
	defer c.Unlock()
	if ping < 0 { // cannot go faster than the connection
		ping = 0
	}
	c.inboundDelay = ping / 2
	c.outboundDelay = ping / 2
}

func (c *connection) GetClientUDPAddress() *net.UDPAddr {
	return c.clientAddress
}

func (c *connection) ForwardServerToClient(ctx context.Context, proxy *net.UDPConn, logger logging.Logger) {
	defer func() {
		logger.Info("closing connection with client %s", c.clientAddress)
		if err := c.close(); err != nil {
			logger.Error(err)
		}
	}()
	packets := make(chan []byte) // unbuffered
	go func() {
		if err := c.readFromServer(packets); err != nil {
			logger.Error(err)
		}
	}()
	for {
		select {
		case packet := <-packets:
			go func() {
				err := writeToClientWithDelay(ctx, c.getOutboundDelay(), proxy, c.clientAddress, packet)
				if err != nil {
					logger.Error(err)
				}
			}()
		case <-ctx.Done():
			logger.Info("context canceled, closing connection")
			c.close()
			return
		}
	}
}

func (c *connection) readFromServer(packets chan<- []byte) error {
	buffer := make([]byte, c.bufferSize)
	for {
		bytesRead, err := c.server.Read(buffer)
		if err != nil {
			return err
		}
		data := make([]byte, bytesRead)
		copy(data, buffer[:bytesRead])
		packets <- data
	}
}

func (c *connection) WriteToServerWithDelay(ctx context.Context, data []byte) error {
	return writeToServerWithDelay(ctx, c.getInboundDelay(), c.server, data)
}

func writeToServerWithDelay(ctx context.Context, delay time.Duration, server UDPConn, data []byte) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	timer := time.AfterFunc(delay, func() {
		err = writeToServer(server, data)
		cancel()
	})
	<-ctx.Done() // done when write is done or context canceled externally
	timer.Stop()
	return err
}

func writeToClientWithDelay(ctx context.Context, delay time.Duration, proxy *net.UDPConn, client *net.UDPAddr, data []byte) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	timer := time.AfterFunc(delay, func() {
		err = writeToClient(proxy, client, data)
		cancel()
	})
	<-ctx.Done() // done when write is done or context is canceled externally
	timer.Stop()
	return err
}

func writeToClient(proxy UDPConn, client *net.UDPAddr, data []byte) error {
	bytesWritten, err := proxy.WriteToUDP(data, client)
	if err != nil {
		return err
	} else if bytesWritten != len(data) {
		return fmt.Errorf("read %d bytes from server and wrote %d bytes to client", len(data), bytesWritten)
	}
	return nil
}

func writeToServer(server UDPConn, data []byte) error {
	n, err := server.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return fmt.Errorf("wrote %d bytes but data was %d bytes", n, len(data))
	}
	return nil
}
