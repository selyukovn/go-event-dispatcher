# Event Dispatcher

### TL;DR

A minimalistic package for working with domain (and other) events.

### Here's the thing

This package was primarily created to avoid copying the domain event dispatcher from one service to another
(in a modular monolith, this code could be in a common module).

---

To accumulate events before processing, a special [`Collection`](events_collection.go) is used.\
This collection can be passed by pointer as an argument to a method that raises events.

```go
evs := event.NewCollection(event.ColOptAsync())
// ...
go func () { a.SomeMethod(arg1, arg2, ..., evs) }()
// ...
b.OtherMethod(arg1, arg2, ..., evs)
// ...
dispatcher.Dispatch(evs)
```

---

[`EventInterface`](event.go) does not require technical methods (like `Name() string`) to link events with handlers.
This clearly makes the domain model cleaner and more independent.

However, often the minimal set of data in an event includes the time when it occurred.
To avoid code duplication, you can use the [`Event`](event.go) type as an embedded field.

---

When using [`DispatcherImplGeneralFn`](dispatcher_impl_general_fn.go),
the event handling flow is defined when creating a dispatcher instance.

The way to determine handlers for a specific event is left to the developer’s discretion.
A suggested simple approach is to use type assertions for events,
which reduces mistakes probability compared to, for example, using strings from `Name() string` methods.

Sync / async and other aspects of events processing are also left to the developer’s discretion.

### Example

```go
// accountInternalEvent -- is an example of inner struct to avoid modifying built-in field "Event" from the outside
// and to contain common fields of all related events.
type accountInternalEvent struct {
    event.Event
    accountId uint
}

type AccountCreatedEvent struct {
    accountInternalEvent
    // ...
}

type AccountDeactivatedEvent struct {
    accountInternalEvent
    // ...
}

func (a *Account) Deactivate(now time.Time, evs *event.Collection) {
    // ...
    evs.Add(NewAccountDeactivatedEvent(a.id))
    // ...
}

// --

eventDispatcher := event.NewDispatcherImplGeneralFn(func (ctx context.Context, e event.EventInterface) {
    switch v := e.(type) {

    // account
    case AccountCreatedEvent:
        // ...
    case AccountDeactivatedEvent:
        // ...

    // unknown event
    default:
        panic(fmt.Errorf("handling of %T event is not defined", v))
    }
})

// --

evs := event.NewCollection()

acc := ...

acc.Deactivate(time.Now(), evs)

eventDispatcher.Dispatch(ctx, evs)
```
