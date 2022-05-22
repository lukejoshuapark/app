package app

import "context"

// Service represents a long-running service that is critical to the operation
// of the program.
//
// The Run method starts the service.  Implementations must exit as soon as
// possible if the provided context ends.
type Service interface {
	Run(ctx context.Context) error
}
