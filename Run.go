package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

// Run runs the provided App.  Should be called directly in main.
func Run(app App) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigc
		cancel()
	}()

	if err := runInContext(ctx, cancel, app); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func runInContext(ctx context.Context, cancel context.CancelFunc, app App) error {
	// Get the services we'll be running.  If there aren't any, we just stop
	// here.
	services, err := app.Services()
	if err != nil {
		return err
	}

	serviceCount := len(services)
	if serviceCount < 1 {
		return nil
	}

	// Start the services and prepare a channel for their error results.
	errc := make(chan error, serviceCount)
	for _, service := range services {
		startService(ctx, errc, service)
	}

	// Wait for our first error, or the context to end.  If we receive an error
	// first, we store it so we can report it as the originating fault, then we
	// cancel the context.  If the context ends first, we're doing a graceful
	// shutdown.
	var terminationError error
	errorCount := 0

	select {
	case <-ctx.Done():
	case err := <-errc:
		terminationError = err
		errorCount++
		cancel()
	}

	// Receive all the errors to ensure all the services have shutdown.
	for errorCount < serviceCount {
		<-errc
		errorCount++
	}

	// Close the error channel and return the termination error, if there was
	// one.
	close(errc)
	return terminationError
}

func startService(ctx context.Context, errc chan<- error, service Service) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errc <- fmt.Errorf("service panic: %v\n\n%v", r, debug.Stack())
			}
		}()

		errc <- service.Run(ctx)
	}()
}
