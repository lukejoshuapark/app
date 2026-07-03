package app

import "context"

// App represents a runtime application.  It can be provided to Run to start an
// application.
type App interface {
	Services(ctx context.Context) ([]Service, error)
}
