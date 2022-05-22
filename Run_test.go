package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/lukejoshuapark/test"
	"github.com/lukejoshuapark/test/is"
)

func SetupRunTest(failsSetup bool, services []Service) (*TestApp, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return NewTestApp(failsSetup, services), ctx, cancel
}

func TestRunExitsWhenSetupFails(t *testing.T) {
	// Arrange.
	app, ctx, cancel := SetupRunTest(true, nil)

	// Act.
	err := runInContext(ctx, cancel, app)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, err.Error(), is.EqualTo("failed setup"))
}

func TestRunExitsWhenNoServices(t *testing.T) {
	// Arrange.
	app, ctx, cancel := SetupRunTest(false, nil)

	// Act.
	err := runInContext(ctx, cancel, app)

	// Assert.
	test.That(t, err, is.Nil)
}

func TestRunExitsCleanlyWhenContextCancelled(t *testing.T) {
	// Arrange.
	services := []Service{NewTimeoutService("Timeout1", 1000*time.Millisecond)}
	app, ctx, cancel := SetupRunTest(false, services)

	// Act.
	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	err := runInContext(ctx, cancel, app)

	// Assert.
	test.That(t, err, is.Nil)
}

func TestRunExitsWhenServiceReturnsError(t *testing.T) {
	// Arrange.
	services := []Service{
		NewTimeoutService("Timeout1", 100*time.Millisecond),
	}

	app, ctx, cancel := SetupRunTest(false, services)

	// Act.
	err := runInContext(ctx, cancel, app)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, err.Error(), is.EqualTo("timeout ended"))
}

func TestRunExitsWhenServicePanics(t *testing.T) {
	// Arrange.
	service := NewTimeoutService("Timeout1", 100*time.Millisecond)
	service.shouldPanic = true

	services := []Service{service}

	app, ctx, cancel := SetupRunTest(false, services)

	// Act.
	err := runInContext(ctx, cancel, app)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, strings.HasPrefix(err.Error(), "service panic: timeout ended"), is.True)
}

// --

type TestApp struct {
	failsSetup bool
	services   []Service
}

var _ App = &TestApp{}

func NewTestApp(failsSetup bool, services []Service) *TestApp {
	return &TestApp{
		failsSetup: failsSetup,
		services:   services,
	}
}

func (a *TestApp) Services() ([]Service, error) {
	if a.failsSetup {
		return nil, errors.New("failed setup")
	}

	return a.services, nil
}

// --

type TimeoutService struct {
	name        string
	shouldPanic bool
	timeout     time.Duration
}

var _ Service = &TimeoutService{}

func NewTimeoutService(name string, timeout time.Duration) *TimeoutService {
	return &TimeoutService{
		name:    name,
		timeout: timeout,
	}
}

func (s *TimeoutService) Run(ctx context.Context) error {
	timer := time.NewTimer(s.timeout)

	select {
	case <-ctx.Done():
		return nil
	case <-timer.C:
		timer.Stop()
		if s.shouldPanic {
			panic("timeout ended")
		}

		return errors.New("timeout ended")
	}
}
