
## [0.1.0] - 2025-11-06

### FEATURES
- Added basic abstractions `DispatcherInterface` and `EventInterface`.
- Added `Dispatcher` as default implementation of the dispatcher interface.\
  Handlers are resolved via user‑supplied function at dispatcher creation. 
- Added `Event` struct with `OccurredAt() time.Time` method to avoid boilerplate in custom events.
- Added `Collection` for accumulating events before processing.\
  Implemented pre-allocation ability via `ColOptInitialSize(int)` and ability to work in async mode via `ColOptAsync()`.
- Basic tests for `Dispatcher` and `Collection`.
