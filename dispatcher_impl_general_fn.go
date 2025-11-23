package event

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type DispatcherImplGeneralFn struct {
	fnHandler func(ctx context.Context, e EventInterface)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDispatcherImplGeneralFn
//
// Panics, if "fnHandler" is nil.
//
// "fnHandler" executes for each event handled by DispatcherImplGeneralFn.Dispatch.
func NewDispatcherImplGeneralFn(fnHandler func(ctx context.Context, e EventInterface)) *DispatcherImplGeneralFn {
	if fnHandler == nil {
		panic(fmt.Errorf("NewDispatcherImplGeneralFn : fnHandler must not be nil"))
	}

	return &DispatcherImplGeneralFn{
		fnHandler: fnHandler,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Dispatch
//
// Panics, if any argument is nil.
func (d *DispatcherImplGeneralFn) Dispatch(ctx context.Context, evs *Collection) {
	if ctx == nil {
		panic(fmt.Errorf("%T.Dispatch : ctx must not be nil", d))
	}

	if evs == nil {
		panic(fmt.Errorf("%T.Dispatch : evs must not be nil", d))
	}

	for _, e := range evs.All() {
		d.fnHandler(ctx, e)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
