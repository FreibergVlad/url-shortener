package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

const ShutdownTimeout = 5 * time.Second

type HttpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type HttpServerWithGracefulShutdown struct {
	HttpServer
	Quit chan os.Signal
}

func NewServer(http HttpServer) *HttpServerWithGracefulShutdown {
	return &HttpServerWithGracefulShutdown{
		HttpServer: http,
		Quit:       make(chan os.Signal, 1),
	}
}

func (s *HttpServerWithGracefulShutdown) Run() {
	panic := make(chan error, 1)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic <- err
		}
	}()

	signal.Notify(s.Quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case signal := <-s.Quit:
		log.Info().Msgf("Shutdown signal '%s' received from OS", signal)
	case err := <-panic:
		log.Panic().Err(err).Msg("Error running server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msgf("Timeout %s exceeded, forcing shutdown", ShutdownTimeout)
	}

	log.Info().Msg("Server gracefully stopped")
}
