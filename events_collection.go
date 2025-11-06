package event

import (
	"fmt"
	"sync"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

// CollectionSizeDefault -- no need to be public, but useful for docs.
//
// DDD-approach recommends to update one aggregate within a single transaction,
// so in most cases only one method of the aggregate fills the collection,
// and number of emitted events is enough small.
const CollectionSizeDefault int = 2
const CollectionSizeMax int = 999_999

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Collection struct {
	evs      []EventInterface
	evsAdded int // to avoid nil-values in result slice, because evs length may be > evsAdded

	async bool
	mu    sync.RWMutex

	initiated bool
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

type ColOpt func(*Collection)

// ColOptAsync
//
// Collection usage will be controlled by sync.RWMutex.
//
// Panics on applying to initialized collection -- see NewCollection().
func ColOptAsync() ColOpt {
	return func(c *Collection) {
		if c.initiated {
			panic(fmt.Errorf("collection option (Async) must not be applied to initialized collection"))
		}

		c.async = true
	}
}

// ColOptInitialSize
//
// Use it, if expected collection size is much more than CollectionSizeDefault.
// If s = 0, CollectionSizeDefault is used instead.
//
// Panics, if not in range [0, CollectionSizeMax].
//
// Panics on applying to initialized collection  -- see NewCollection().
func ColOptInitialSize(s int) ColOpt {
	return func(c *Collection) {
		if c.initiated {
			panic(fmt.Errorf("collection option (InitialSize) must not be applied to initialized collection"))
		}

		if s == 0 {
			s = CollectionSizeDefault
		}

		if !(0 <= s && s <= CollectionSizeMax) {
			panic(fmt.Errorf(
				"collection option (InitialSize) expects size in range [0, %d], but %d given",
				CollectionSizeMax,
				s,
			))
		}

		c.evs = make([]EventInterface, 0, s)
	}
}

// --

// NewCollection
//
// Panics, if opts contains nil.
func NewCollection(opts ...ColOpt) *Collection {
	c := &Collection{
		evs:       nil,
		evsAdded:  0,
		async:     false,
		mu:        sync.RWMutex{},
		initiated: false,
	}

	// options
	for _, opt := range opts {
		if opt == nil {
			panic(fmt.Errorf("NewCollection cannot handle nil option"))
		}

		opt(c)
	}

	// missed options defaults
	if c.evs == nil {
		c.evs = make([]EventInterface, 0, CollectionSizeDefault)
	}

	// finish
	c.initiated = true

	return c
}

func (c *Collection) assertUsage() {
	if !c.initiated {
		panic(fmt.Errorf("unable to use %T before initialization", c))
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Panics, if argument e is nil.
// Panics, if too many events are added -- see CollectionSizeMax.
// Panics, if collection is not initiated -- see NewCollection.
func (c *Collection) Add(e EventInterface) {
	if e == nil {
		panic(fmt.Errorf("%T cannot collect nils", c))
	}

	c.assertUsage()

	if c.async {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if c.evsAdded == CollectionSizeMax {
		panic(fmt.Errorf("%T: has too many events (%d)", c, CollectionSizeMax))
	}

	if c.evsAdded == len(c.evs) {
		c.evs = append(c.evs, e)
	} else {
		c.evs[c.evsAdded-1] = e
	}

	c.evsAdded++
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// Len
//
// Panics, if collection is not initiated -- see NewCollection.
func (c *Collection) Len() int {
	c.assertUsage()

	if c.async {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	return c.evsAdded
}

// IsEmpty
//
// Panics, if collection is not initiated -- see NewCollection.
func (c *Collection) IsEmpty() bool {
	c.assertUsage()

	if c.async {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	return c.evsAdded == 0
}

// All
//
// Returns a COPY of the collected events to avoid modifying internal state from the outside.
// It is not expected to collect massive sets of events or use this method multiple times,
// so performance issues are not expected accordingly.
//
// Panics, if collection is not initiated -- see NewCollection.
func (c *Collection) All() []EventInterface {
	c.assertUsage()

	if c.async {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	cloned := make([]EventInterface, c.evsAdded)
	copy(cloned, c.evs[:c.evsAdded])
	return cloned
}
