package event

import (
	"context"
)

type DispatcherInterface interface {
	// Dispatch
	//
	// Panics, if any argument is nil.
	Dispatch(ctx context.Context, evs *Collection)
}
