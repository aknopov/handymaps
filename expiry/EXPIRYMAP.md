# ExpiryMap User Guide

## Map Creation

A map is created with a call to `NewExpiryMap`. For example - `exMap := expiry.NewExpiryMap[int, string]()`. The first generic argument is the key type, and the second is the value type.
The map created with this call does not expire entries and has unlimited capacity. You need to set the capacity and expiry time to reach full functionality.

## Method Chaining

Methods that configure `ExpiryMap` can be chained in any order.

## Loader Function

A user-defined loader function is invoked synchronously on the first `Get` call for a key. Subsequent calls perform a non-blocking read-through operation until the key expires.
The loading function should have the signature `func(key K) (V, error)`. For example - `exMap.WithLoader(func(key string) (int, error) { return len(key), nil })`.
If the load fails, the function should return an error that is returned as the second value of the `Get` call. The first value is the "zero" value of type `K`.

## Thread Safety and Blocking

All major cache operations are thread-safe and use a Read-Write locking mechanism. Operations such as `Capacity`, `ExpireTime`, `Len`, and `Peek` either do not block or allow multiple read operations.
Operations `Clear` and `Discard` block all operations with a Write lock. The `Get` operation starts with a Read lock that might be upgraded to a Write lock if the value needs to be loaded.

## Life Cycle

The creation of an expiry map starts a goroutine that evicts expired keys and listens to a particular channel for timer events. `ExpiryMap` provides a `Discard` method that stops this goroutine.
This marks the end of the map's life, and most operations will result in a panic thereafter.

## Listening to Cache Events

`ExpiryMap` allows tracking map events that could be used, for example, in collecting statistics. The map allows unlimited `Listener` instances that can be added with `AddListener` and removed with `RemoveListener`. These listeners are invoked synchronously for each event in the order of their addition. The map provides the following events: adding, expiring, peeking, removing, missing (`Peek` operation), replacing, and load failures.
