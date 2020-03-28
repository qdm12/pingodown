package main

import (
	"context"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/pingodown/internal/proxy"
)

func main() {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, 0)
	if err != nil {
		panic(err)
	}

	envParams := params.NewEnvParams()
	listenAddresss, err := envParams.GetEnv("LISTEN_ADDRESS", params.Default(":8000"))
	if err != nil {
		logger.Error(err)
		return
	}
	serverAddress, err := envParams.GetEnv("SERVER_ADDRESS", params.Compulsory())
	if err != nil {
		logger.Error(err)
		return
	}
	ping, err := envParams.GetDuration("PING", params.Default("100ms"))
	if err != nil {
		logger.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	proxy, err := proxy.NewProxy(listenAddresss, serverAddress, logger, ping)
	if err != nil {
		logger.Error(err)
		return
	}
	if err := proxy.Run(ctx); err != nil {
		logger.Error(err)
	}
}
