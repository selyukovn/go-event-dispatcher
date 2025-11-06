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

When using [`Dispatcher`](dispatcher.go),
the set of event handlers is defined when creating a dispatcher instance (for now).\
This allows using either a single dispatcher instance for all events
or multiple instances for some groups of events or so -- whatever.\

The way to determine handlers for a specific event is left to the developer’s discretion. 
A recommended simple approach is to use type assertions for events, 
which reduces the probability of errors compared to, for example, using strings from `Name() string` methods.

Sync / async and other aspects of events processing are also left to the developer’s discretion.

### Example

```go
// package-only struct example to avoid modifying built-in field Event from the outside
// and to contain common fields of all related events.
type accountInternalEvent struct {
    event.Event
    accountId uint
}

type AccountActivatedEvent struct {
    accountInternalEvent
    // ...
}

// --

eventDispatcher := func () event.DispatcherInterface {
    emptySet := []func (ctx context.Context, e event.EventInterface){}
    
    return event.NewDispatcher(func (e event.EventInterface) []func (ctx context.Context, e event.EventInterface) {
        switch v := e.(type) {

        // account
        case AccountActivatedEvent:
            return []func (ctx context.Context, e event.EventInterface){
                // handler 1 
                func (ctx context.Context, e event.EventInterface) {
                    // ...
                },

                // handler 2
                // ...
            }

        // ...

        // not registered
        default:
            panic(fmt.Errorf("there are no handlers registered for %T event", v))
        }

        return emptySet
    })
}()

// --

evs := event.NewCollection()

_ = account.Activate(time.Now(), evs) // e.g. evs.Add(NewAccountActivatedEvent(a.id))

eventDispatcher.Dispatch(ctx, evs)
```
