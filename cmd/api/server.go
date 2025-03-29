package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.server.port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	serverShutdownErr := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		quitSignal := <-quit
		signal.Stop(quit)

		log.Printf("quit signal: %q received. starting graceful shutdown\n", quitSignal.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			serverShutdownErr <- err
			return
		}

		serverShutdownErr <- nil
	}()

	app.logger.Printf("server running on port %d\n", app.cfg.server.port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-serverShutdownErr; err != nil {
		return err
	}

	return nil
}
