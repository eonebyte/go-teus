package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func Run(
	app Application,
	timeout time.Duration,
	cleanup func() error,
) error {

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Start()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {

	case err := <-serverErrors:

		if err != nil &&
			!errors.Is(err, context.Canceled) {
			return err
		}

		return nil

	case <-quit:

		ctx, cancel := context.WithTimeout(
			context.Background(),
			timeout,
		)
		defer cancel()

		if err := app.Shutdown(ctx); err != nil {
			return fmt.Errorf(
				"shutdown failed: %w",
				err,
			)
		}

		if cleanup != nil {
			if err := cleanup(); err != nil {
				return fmt.Errorf(
					"cleanup failed: %w",
					err,
				)
			}
		}

		return nil
	}
}