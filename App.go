package app

// App represents a runtime application.  It can be provided to Run to start an
// application.
type App interface {
	Services() ([]Service, error)
}
