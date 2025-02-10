package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const ShutdownTimeout = 5 * time.Second

type ServerWithGracefulShutdown struct {
	*grpc.Server
	listener net.Listener
	Quit     chan os.Signal
}

func NewServer(listener net.Listener, opt ...grpc.ServerOption) *ServerWithGracefulShutdown {
	return &ServerWithGracefulShutdown{
		grpc.NewServer(opt...),
		listener,
		make(chan os.Signal, 1),
	}
}

func (s *ServerWithGracefulShutdown) Run() {
	panicChan := make(chan error, 1)

	go func() {
		if err := s.Serve(s.listener); err != nil {
			panicChan <- err
		}
	}()

	signal.Notify(s.Quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case signal := <-s.Quit:
		log.Info().Msgf("Shutdown signal '%s' received from OS", signal)
	case err := <-panicChan:
		log.Panic().Err(err).Msg("Error running server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	stop := make(chan struct{}, 1)
	go func() {
		s.GracefulStop()
		stop <- struct{}{}
	}()

	select {
	case <-stop:
		log.Info().Msg("Server gracefully stopped")
	case <-ctx.Done():
		log.Error().Msgf("Timeout %s exceeding, forcing shutdown", ShutdownTimeout)
		s.Stop()
	}
}
