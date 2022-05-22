![](icon.png)

# app

Easy orchestration and management of services at runtime!

## Usage Example

```go
package main

import (
	"fmt"

	"github.com/lukejoshuapark/app"
)

func main() {
	app.Run(&App{})
}

// --

type App struct{}

func (a *App) Services() ([]app.Service, error) {
	return []app.Service{
		NewWaitingService(),
	}, nil
}

// --

type WaitingService struct{}

var _ app.Service = &WaitingService{}

func NewWaitingService() *WaitingService {
	return &WaitingService{}
}

func Run(ctx context.Context) error {
	fmt.Println("I'm just going to sit here until the application terminates...")
	<-ctx.Done()

	fmt.Println("Goodbye!")
	return nil
}
```

## Behavior

This module provides a lightweight but powerful runtime service management
layer.  It has two concepts:

- `App` - A type that implements the `Services` method.  It simply constructs
and returns all the runtime services that make up the application.

- `Service` An actual runtime service.  A service can do whatever it likes, as
long as it:

	- Returns a nil error as soon as possible when the context provided to it
	ends.

	- Returns an error when something unrecoverable occurs.

Passing your implementation of `app.App` to `app.Run` starts the application.
The runtime layer handles termination signals by cancelling the context provided
to each service.

---

Icons made by [justicon](https://www.flaticon.com/authors/justicon) from
[www.flaticon.com](https://www.flaticon.com/).
