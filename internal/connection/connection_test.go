package connection

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/pingodown/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_close(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := mocks.NewMockUDPConn(mockCtrl)
	mockUDPConn.EXPECT().Close().Return(fmt.Errorf("error")).Times(1)
	c := &connection{
		server: mockUDPConn,
	}
	err := c.close()
	require.Error(t, err)
	assert.Equal(t, "error", err.Error())
}

func Test_getInboundDelay(t *testing.T) {
	t.Parallel()
	c := &connection{inboundDelay: time.Second}
	assert.Equal(t, time.Second, c.getInboundDelay())
}

func Test_getOutboundDelay(t *testing.T) {
	t.Parallel()
	c := &connection{outboundDelay: time.Second}
	assert.Equal(t, time.Second, c.getOutboundDelay())
}

func Test_SetPing(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		ping          time.Duration
		inboundDelay  time.Duration
		outboundDelay time.Duration
	}{
		"zero ping": {},
		"negative ping": {
			ping: -time.Second,
		},
		"positive ping": {
			ping:          100 * time.Millisecond,
			inboundDelay:  50 * time.Millisecond,
			outboundDelay: 50 * time.Millisecond,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := &connection{}
			c.SetPing(tc.ping)
			assert.Equal(t, tc.inboundDelay, c.inboundDelay)
			assert.Equal(t, tc.outboundDelay, c.outboundDelay)
		})
	}
}

func Test_PingApproximatesTo(t *testing.T) {
	t.Parallel()
	const inboundDelay, outboundDelay = 100 * time.Millisecond, 100 * time.Millisecond
	tests := map[string]struct {
		d            time.Duration
		approximates bool
	}{
		"same ping": {
			d:            200 * time.Millisecond,
			approximates: true,
		},
		"minus 3ms": {
			d:            197 * time.Millisecond,
			approximates: true,
		},
		"plus 3ms": {
			d:            203 * time.Millisecond,
			approximates: true,
		},
		"plus 10ms": {
			d:            210 * time.Millisecond,
			approximates: false,
		},
		"minus 10ms": {
			d:            190 * time.Millisecond,
			approximates: false,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := &connection{
				inboundDelay:  inboundDelay,
				outboundDelay: outboundDelay,
			}
			approximates := c.PingApproximatesTo(tc.d)
			assert.Equal(t, tc.approximates, approximates)
		})
	}
}

func Test_GetClientUDPAddress(t *testing.T) {
	t.Parallel()
	udpAddress := &net.UDPAddr{
		IP:   net.IP{10, 10, 10, 10},
		Port: 8000,
	}
	c := &connection{clientAddress: udpAddress}
	assert.Equal(t, udpAddress, c.GetClientUDPAddress())
}
