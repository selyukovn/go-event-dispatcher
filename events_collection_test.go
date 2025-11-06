package event

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// NewCollection
// ---------------------------------------------------------------------------------------------------------------------

func Test_NewCollection(t *testing.T) {
	// Bad Client
	// --------------------------------

	t.Run("PanicBadClient", func(t *testing.T) {
		assert.Panics(t, func() { NewCollection(nil) })
		assert.Panics(t, func() { NewCollection(ColOptAsync(), nil) })
		assert.Panics(t, func() { NewCollection(nil, ColOptInitialSize(5)) })
		assert.Panics(t, func() { NewCollection(ColOptAsync(), nil, ColOptInitialSize(5)) })
	})

	// Empty
	// --------------------------------

	t.Run("Empty", func(t *testing.T) {
		for _, evs := range []*Collection{
			NewCollection(),
			NewCollection(ColOptAsync()),
			NewCollection(ColOptInitialSize(5)),
			NewCollection(ColOptAsync(), ColOptInitialSize(5)),
		} {
			assert.NotNil(t, evs)
			assert.Zero(t, evs.Len())
			assert.True(t, evs.IsEmpty())
			assert.Empty(t, evs.All())
		}
	})

	// Collection Option - Async
	// --------------------------------

	t.Run("ColOptAsync_BadClient", func(t *testing.T) {
		assert.Panics(t, func() {
			ColOptAsync()(NewCollection())
		})
	})

	t.Run("ColOptAsync_AsyncAdd", func(t *testing.T) {
		evs := NewCollection(ColOptAsync())

		type TE struct {
			w int
			j int
		}

		const numWorkers = 5
		const eventsPerWorker = 10
		var wg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(w int) {
				defer wg.Done()
				for j := 0; j < eventsPerWorker; j++ {
					evs.Add(TE{w, j})
				}
				time.Sleep(100 * time.Millisecond)
			}(w)
		}
		wg.Wait()

		assert.False(t, evs.IsEmpty())
		assert.Equal(t, evs.Len(), len(evs.All()))
		assert.Equal(t, numWorkers*eventsPerWorker, evs.Len())
		m := make(map[int]map[int]bool)
		for w := 0; w < numWorkers; w++ {
			m[w] = make(map[int]bool)
			for j := 0; j < eventsPerWorker; j++ {
				m[w][j] = false
			}
		}
		for _, e := range evs.All() {
			assert.False(t, m[e.(TE).w][e.(TE).j], "duplicated events")
			m[e.(TE).w][e.(TE).j] = true
		}
		for w := 0; w < numWorkers; w++ {
			for j := 0; j < eventsPerWorker; j++ {
				assert.True(t, m[w][j], "not all events added")
			}
		}
	})

	// Collection Option - InitialSize
	// --------------------------------

	t.Run("ColOptInitialSize_BadClient", func(t *testing.T) {
		// on initiated collection
		assert.Panics(t, func() {
			ColOptInitialSize(5)(NewCollection())
		})

		// min size
		assert.Panics(t, func() {
			NewCollection(ColOptInitialSize(-1))
		})

		// max size
		assert.Panics(t, func() {
			NewCollection(ColOptInitialSize(CollectionSizeMax + 1))
		})
	})

	t.Run("ColOptInitialSize", func(t *testing.T) {
		assert.NotPanics(t, func() {
			evs := NewCollection(ColOptInitialSize(0))
			evs.Add(struct{}{})
		})
		assert.NotPanics(t, func() {
			evs := NewCollection(ColOptInitialSize(CollectionSizeMax))
			evs.Add(struct{}{})
		})
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Add
// ---------------------------------------------------------------------------------------------------------------------

func Test_Collection_Add(t *testing.T) {
	type TE struct{ i int }

	// Bad client
	// --------------------------------

	t.Run("PanicBadClient", func(t *testing.T) {
		// not initiated
		evs := &Collection{}
		assert.Panics(t, func() {
			evs.Add(TE{1})
		})

		// nil event
		evs = NewCollection()
		assert.Panics(t, func() {
			evs.Add(nil)
		})

		// max elements
		evs = NewCollection(ColOptInitialSize(CollectionSizeMax))
		var i int
		for i = 0; i < CollectionSizeMax; i++ {
			evs.Add(struct{}{})
		}
		assert.Panics(t, func() {
			evs.Add(struct{}{})
		})
	})

	// Success
	// --------------------------------

	t.Run("Success", func(t *testing.T) {
		evs := NewCollection()

		const n = 5

		for i := 0; i < n; i++ {
			evs.Add(TE{i})
		}

		assert.False(t, evs.IsEmpty())
		assert.Equal(t, n, evs.Len())
		assert.Equal(t, evs.Len(), len(evs.All()))

		for i, ev := range evs.All() {
			assert.Equal(t, i, ev.(TE).i)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// All
// ---------------------------------------------------------------------------------------------------------------------

func TestCollection_All(t *testing.T) {
	t.Run("PanicBadClient", func(t *testing.T) {
		evs := &Collection{}

		assert.Panics(t, func() {
			evs.All()
		})
	})

	t.Run("Copy", func(t *testing.T) {
		evs := NewCollection()
		evs.Add(struct{}{})
		evs.Add(struct{}{})

		cp := evs.All()
		cp = append(cp, struct{}{})
		cp = append(cp, struct{}{})
		cp = append(cp, struct{}{})
		assert.Len(t, cp, 5)

		assert.Len(t, evs.All(), 2)
	})
}
