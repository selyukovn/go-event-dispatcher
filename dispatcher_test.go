package event

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
// NewDispatcher
// ---------------------------------------------------------------------------------------------------------------------

func Test_NewDispatcher(t *testing.T) {
	t.Run("PanicForBadClient", func(t *testing.T) {
		// nil constructor
		assert.Panics(t, func() {
			NewDispatcher(nil)
		})

		// nil handlers set
		assert.Panics(t, func() {
			ctx := context.TODO()
			evs := NewCollection()
			evs.Add(struct{}{})

			NewDispatcher(func(e EventInterface) []func(ctx context.Context, e EventInterface) {
				return nil
			}).Dispatch(ctx, evs)
		})

		// nil handler in set
		assert.Panics(t, func() {
			ctx := context.TODO()
			evs := NewCollection()
			evs.Add(struct{}{})

			NewDispatcher(func(e EventInterface) []func(ctx context.Context, e EventInterface) {
				return []func(ctx context.Context, e EventInterface){
					nil,
				}
			}).Dispatch(ctx, evs)
		})
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Dispatch
// ---------------------------------------------------------------------------------------------------------------------

func Test_Dispatcher_Dispatch(t *testing.T) {
	// Success
	// --------------------------------

	t.Run("SuccessEmptyCollection", func(t *testing.T) {
		ctx := context.TODO()
		evs := NewCollection()
		d := NewDispatcher(func(e EventInterface) []func(ctx context.Context, e EventInterface) {
			return []func(ctx context.Context, e EventInterface){
				func(ctx context.Context, e EventInterface) {
					panic("should not happen")
				},
			}
		})

		assert.NotPanics(t, func() {
			d.Dispatch(ctx, evs)
		})
	})

	t.Run("SuccessAllHandlers", func(t *testing.T) {
		type TE struct {
			handler1 uint
			handler2 uint
		}
		tes := []*TE{{}, {}, {}}

		ctxTestKey := "test"
		ctxTestValue := 123
		ctx := context.TODO()
		ctx = context.WithValue(ctx, ctxTestKey, ctxTestValue)

		evs := NewCollection()

		for _, te := range tes {
			evs.Add(te)
		}

		d := NewDispatcher(func(e EventInterface) []func(ctx context.Context, e EventInterface) {
			return []func(ctx context.Context, e EventInterface){
				func(ctx context.Context, e EventInterface) {
					te := e.(*TE)
					te.handler1++

					// check same context
					assert.Equal(t, ctx.Value(ctxTestKey), ctxTestValue)
				},
				func(ctx context.Context, e EventInterface) {
					te := e.(*TE)
					te.handler2++

					// check same context
					assert.Equal(t, ctx.Value(ctxTestKey), ctxTestValue)
				},
			}
		})

		d.Dispatch(ctx, evs)

		for _, te := range tes {
			assert.Equal(t, te.handler1, uint(1))
			assert.Equal(t, te.handler2, uint(1))
		}
	})

	// Bad client
	// --------------------------------

	t.Run("PanicBadClient", func(t *testing.T) {
		ctx := context.TODO()
		evs := NewCollection()
		d := NewDispatcher(func(e EventInterface) []func(ctx context.Context, e EventInterface) {
			return []func(ctx context.Context, e EventInterface){
				func(ctx context.Context, e EventInterface) {
					// ...
				},
			}
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
}
