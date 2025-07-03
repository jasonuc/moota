package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasonuc/moota/internal/events"
)

func (app *application) serve() error {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.server.port),
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  1 * time.Minute,
		ErrorLog:     app.logger,
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
	go func() {
		err := app.routers.EventsRouter.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := app.routers.SSERouter.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		// This goroutine simulates some events being published in the background
		ctx := context.Background()
		for {
			if rand.Intn(2) == 0 {
				_ = app.routers.EventBus.Publish(ctx, events.SeedPlanted{})
			} else {
				_ = app.routers.EventBus.Publish(ctx, events.SeedGenerated{})
			}

			time.Sleep(time.Millisecond * time.Duration(3000+rand.Intn(5000)))
		}
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
