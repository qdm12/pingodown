package proxy

import (
	"context"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/pingodown/internal/ping"
	"github.com/qdm12/pingodown/internal/state"
)

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

func pingClients(pinger ping.Pinger, state state.State, logger logging.Logger) {
	for _, address := range state.GetClientAddresses() {
		address := address
		latency, err := pinger.GetLatency(address)
		if err != nil {
			logger.Error(err)
			continue
		}
		state.SetLatency(address, latency)
	}
}

func updatePings(state state.State, logger logging.Logger) {
	highestLatency := state.GetHighestLatency()
	for _, address := range state.GetClientAddresses() {
		address := address
		conn, err := state.GetConnection(address)
		if err != nil {
			logger.Error(err)
			continue
		}
		currentLatency, err := state.GetLatency(address)
		if err != nil {
			logger.Error(err)
			continue
		}
		ping := highestLatency - currentLatency
		if !conn.PingApproximatesTo(ping) {
			logger.Info("Setting new ping of %s to %dms", conn.GetClientUDPAddress().String(), ping.Milliseconds())
			conn.SetPing(ping) // works by pointer
		}
	}
}
