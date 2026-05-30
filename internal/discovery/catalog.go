package discovery

import (
	"errors"
	"sync"

	"github.com/mitslab-de/omniinstall/internal/app"
)

// Catalog is an in-memory registry of Application metadata.
// It is safe for concurrent use.
type Catalog struct {
	mu    sync.RWMutex
	byID  map[string]*app.Application
	items []*app.Application
}

// NewCatalog creates an empty Catalog.
func NewCatalog() *Catalog {
	return &Catalog{
		byID: make(map[string]*app.Application),
	}
}

// Add validates and registers an application in the catalog.
// Returns an error if the application fails validation or its ID is already
// registered.
func (c *Catalog) Add(a *app.Application) error {
	if a == nil {
		return errors.New("application must not be nil")
	}
	if err := a.Validate(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.byID[a.ID]; exists {
		return errors.New("duplicate application id: " + a.ID)
	}
	c.byID[a.ID] = a
	c.items = append(c.items, a)
	return nil
}

// Get returns the application with the given ID, or nil if not found.
func (c *Catalog) Get(id string) *app.Application {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.byID[id]
}

// All returns a snapshot of all registered applications in insertion order.
func (c *Catalog) All() []*app.Application {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]*app.Application, len(c.items))
	copy(result, c.items)
	return result
}

// Size returns the number of applications in the catalog.
func (c *Catalog) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
