package event

import "time"

// #####################################################################################################################
// INTERFACE
// #####################################################################################################################

type EventInterface interface{}

// #####################################################################################################################
// DEFAULT IMPLEMENTATION
// #####################################################################################################################

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

// Event
//
// Logically any event is occurred at some point of time,
// so Event can be used as built-in field to simplify your code
// by avoiding of duplicated OccurredAt() method and field in your specific events.
//
// E.g.
//
//	type userInternalEvent struct {
//		event.Event
//		userId uint
//	}
//
//	type UserCreatedEvent struct {
//		userInternalEvent
//		// ...
//	}
//
//	func NewUserCreatedEvent(occurredAt time.Time, userId uint) UserCreatedEvent {
//		return UserCreatedEvent{
//			userInternalEvent: userInternalEvent{Event: event.NewEvent(occurredAt), userId: userId},
//			// ...
//		}
//	}
type Event struct {
	occurredAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEvent
//
// See Event
func NewEvent(occurredAt time.Time) Event {
	return Event{occurredAt: occurredAt}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e Event) OccurredAt() time.Time {
	return e.occurredAt
}

// #####################################################################################################################
