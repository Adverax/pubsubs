# pubsubs

## Overview
`pubsubs` is a Go package that provides a simple and flexible publish-subscribe (pub/sub) pattern implementation. It allows you to create subjects, subscribe to them, and publish events to all subscribers.

## Features
- Generic support for any event type
- Context-aware execution
- Synchronous and asynchronous event handling
- Graceful shutdown of subscribers

## Installation
To install the package, run:
```sh
go get github.com/yourusername/pubsubs
```

# Usage
Define an Event
Define your event type:

```go
type MyEvent struct {
    Message string
}
```

Create a PubSub
Create a new PubSub instance:

```go
ctx := context.Background()
ps, err := NewBuilder[*myEvent]().
    Subject("test").
    Build()
if err != nil {
    log.Fatal(err)
}
defer ps.Close(ctx)
```
Subscribe to Events
Subscribe to events with a handler:

```go
ev := &myEvent{name: "event"}
err = ps.Publish(ctx, ev).WaitWithContext(ctx)
if err != nil {
    log.Fatal(err)
}
```

Publish Events
Publish events to the PubSub:
```go
ev := &myEvent{name: "event"}
err = ps.Publish(ctx, ev).WaitWithContext(ctx)
if err != nil {
    log.Fatal(err)
}
```

Testing
To run tests, use:
```sh
go test ./...
```
