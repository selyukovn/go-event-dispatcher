package event

import (
	"context"
	"fmt"
)

// #####################################################################################################################
// INTERFACE
// #####################################################################################################################

type DispatcherInterface interface {
	// Dispatch
	//
	// Panics, if any argument is nil.
	Dispatch(ctx context.Context, evs *Collection)
}

// #####################################################################################################################
// DEFAULT IMPLEMENTATION
// #####################################################################################################################

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Dispatcher struct {
	fnProvideHandlers func(e EventInterface) []func(ctx context.Context, e EventInterface)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDispatcher
//
// Panics, if any argument is nil.
//
// Dispatcher.Dispatch panics, if fnProvideHandlers returns nil or any element from the returned slice is nil.
func NewDispatcher(
	fnProvideHandlers func(e EventInterface) []func(ctx context.Context, e EventInterface),
) *Dispatcher {
	if fnProvideHandlers == nil {
		panic(fmt.Errorf("NewDispatcher : fnProvideHandlers must not be nil"))
	}

	return &Dispatcher{
		fnProvideHandlers: fnProvideHandlers,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d *Dispatcher) dispatchEvent(ctx context.Context, e EventInterface) {
	fnHandlers := d.fnProvideHandlers(e)

	if fnHandlers == nil {
		panic(fmt.Errorf("%T : fnProvideHandlers provided nil instaead of slice for %#v event", d, e))
	}

	for _, fnHandler := range fnHandlers {
		if fnHandler == nil {
			panic(fmt.Errorf("%T : provided handler cannot be nil for %#v event", d, e))
		}

		fnHandler(ctx, e)
	}
}

// Dispatch
//
// Panics, if any argument is nil.
//
// Panics in case of incorrect fnProvideHandlers -- see NewDispatcher().
func (d *Dispatcher) Dispatch(ctx context.Context, evs *Collection) {
	if ctx == nil {
		panic(fmt.Errorf("%T.Dispatch : ctx must not be nil", d))
	}

	if evs == nil {
		panic(fmt.Errorf("%T.Dispatch : evs must not be nil", d))
	}

	for _, e := range evs.All() {
		d.dispatchEvent(ctx, e)
	}
}

// #####################################################################################################################
