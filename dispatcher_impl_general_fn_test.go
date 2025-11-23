package event

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
// NewDispatcherImplGeneralFn
// ---------------------------------------------------------------------------------------------------------------------

func Test_NewDispatcherImplGeneralFn(t *testing.T) {
	t.Run("PanicForBadClient", func(t *testing.T) {
		assert.Panics(t, func() {
			NewDispatcherImplGeneralFn(nil)
		})

		assert.NotPanics(t, func() {
			NewDispatcherImplGeneralFn(func(ctx context.Context, e EventInterface) {})
		})
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Dispatch
// ---------------------------------------------------------------------------------------------------------------------

func Test_DispatcherImplGeneralFn_Dispatch(t *testing.T) {
	// Bad client
	// --------------------------------

	t.Run("PanicBadClient", func(t *testing.T) {
		ctx := context.TODO()
		evs := NewCollection()
		d := NewDispatcherImplGeneralFn(func(ctx context.Context, e EventInterface) {
			// ...
		})

		// nil context
		assert.Panics(t, func() {
			d.Dispatch(nil, evs)
		})

		// nil collection
		assert.Panics(t, func() {
			d.Dispatch(ctx, nil)
		})
	})

	// Success
	// --------------------------------

	t.Run("SuccessEmptyCollection", func(t *testing.T) {
		ctx := context.TODO()
		evs := NewCollection()
		d := NewDispatcherImplGeneralFn(func(ctx context.Context, e EventInterface) {
			panic("should not happen")
		})
		assert.NotPanics(t, func() {
			d.Dispatch(ctx, evs)
		})
	})

	t.Run("SuccessAll", func(t *testing.T) {
		ctx := context.TODO()
		ctxTestKey := "test"
		ctxTestValue := 123
		ctx = context.WithValue(ctx, ctxTestKey, ctxTestValue)

		type E struct {
			id      int
			handled bool
		}
		e1 := &E{1, false}
		e2 := &E{2, false}
		e3 := &E{3, false}

		evs := NewCollection()
		evs.Add(e1)
		evs.Add(e2)
		evs.Add(e3)

		d := NewDispatcherImplGeneralFn(func(ctx context.Context, e EventInterface) {
			e.(*E).handled = true

			// check same context
			assert.Equal(t, ctx.Value(ctxTestKey), ctxTestValue)
		})

		d.Dispatch(ctx, evs)

		assert.True(t, e1.handled)
		assert.True(t, e2.handled)
		assert.True(t, e3.handled)
	})
}
