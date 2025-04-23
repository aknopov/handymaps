# ExpiryMap User Guide

## Map Creation

A map is created with a call to `NewExpiryMap`. For example - 
```go
exMap := expiry.NewExpiryMap[int, string]()
```
The first generic argument is the key type, and the second is the value type.
The map created with this call does not expire entries and has unlimited capacity. You need to set the capacity and expiry time to have full functionality.

## Method Chaining

Methods that configure `ExpiryMap` can be chained in any order, like
```go
expiryMap := expiry.NewExpiryMap[string, int]().
    WithLoader(func(key string) (int, error) { return len(key), nil }).
    WithMaxCapacity(2).
    ExpireAfter(50 * time.Millisecond)
```

## Loader Function

A user-defined loader function is invoked synchronously on the first `Get` call for a key. Subsequent calls perform a non-blocking read-through operation until the key expires.
The loading function should have the signature `func(key K) (V, error)`. If the load fails, the function should return an error that is returned as the second value of the `Get` call.
The first value in this case is the "zero" value of the type `V`.

## Thread Safety and Blocking

All major cache operations are thread-safe and use a Read-Write locking mechanism. Operations such as `Capacity`, `ExpireTime`, `Len`, and `Peek` either do not block or allow multiple read operations.
In contrary, operations `Clear` and `Discard` block all operations with a write lock. The `Get` operation is different from others.
 It starts with a read lock that might be upgraded to a write lock if the value needs to be loaded.

## Life Cycle

Upon creation of an expiry map, its code starts a goroutine that evicts expired keys and listens to a particular channel for timer events in an infinite loop.
`ExpiryMap` provides a `Discard` method that removes all entries and stops this goroutine. This marks the end of the map's life, and most operations will result in a panic thereafter.

## Listeners

`ExpiryMap` allows tracking map events that could be used, for example, in collecting statistics. The map allows unlimited `Listener` instances that can be added with `AddListener` and removed with `RemoveListener` calls. These listeners are invoked synchronously on each event in the order of their insertion. The map provides the following events: adding, expiring, peeking, removing, missing (`Peek` operation), replacing, and load failures.
