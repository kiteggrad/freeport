package freeport

import (
	"errors"
	"fmt"
	"sync"
)

// global Generator for pkg methods (freeport.Get, freeport.MustGet, ...)
var globGenerator = NewGenerator()

type Generator struct {
	used  map[int]struct{}
	mutex sync.Mutex
}

func NewGenerator() *Generator {
	return &Generator{
		used:  map[int]struct{}{},
		mutex: sync.Mutex{},
	}
}

// Get - returns free port that wasn't used before (from this generator).
func (g *Generator) Get() (port int, err error) {
	// max retries count to get unused port
	const maxRetryCount = 10

	g.mutex.Lock()
	defer g.mutex.Unlock()

	for i := 0; i < maxRetryCount; i++ {
		port, err = tryGetFreeport()
		if err != nil {
			return 0, fmt.Errorf("failed to tryGetFreeport: %w", err)
		}

		// aviod using port that was already requested (but not served)
		if _, used := g.used[port]; !used {
			g.used[port] = struct{}{}
			return port, nil
		}
	}

	return 0, errors.New("failed to find unused port")
}

// MustGet - returns free port that wasn't used before (from this generator) or panics.
func (g *Generator) MustGet() (port int) {
	port, err := g.Get()
	if err != nil {
		panic(fmt.Errorf("failed to GetFreeport: %w", err))
	}
	return port
}
