# Day 03 — Goroutine Syntax

## Spawning a goroutine
```go
go someFunction()       // run someFunction concurrently
go func() {            // anonymous goroutine — inline function
    // work here
}()                    // () at the end calls it immediately
```
Plain English: `go` is the entire keyword. It says "run this in the background, don't wait for it."

## runtime.NumGoroutine
```go
n := runtime.NumGoroutine()   // returns current number of live goroutines
                              // includes main goroutine + test runner goroutine
```

## context.Background and WithCancel
```go
ctx := context.Background()              // root context — never cancelled
ctx, cancel := context.WithCancel(ctx)   // derived context with a cancel function
cancel()                                 // fires cancellation signal to all goroutines using ctx
defer cancel()                           // idiomatic — always defer cancel to avoid leak
```
Plain English: `ctx` is a signal carrier. `cancel()` is a broadcast: "everyone using this ctx, please stop."

## ctx.Done channel
```go
<-ctx.Done()    // blocks until cancel() is called
                // returns immediately once cancel fires
                // type is <-chan struct{} — receive-only, carries no value
```

## Select with context exit
```go
select {
case v := <-ch:         // received a value — do work
    process(v)
case <-ctx.Done():      // cancellation signal — exit cleanly
    return
}
```
Plain English: "Wait for either a value OR a cancel signal — whichever comes first."

## Goroutine exit rule
```go
// Every goroutine must have a clear answer to: "how does this exit?"
go func() {
    for {
        select {
        case v := <-work:
            handle(v)
        case <-ctx.Done():   // exit path — ALWAYS include this
            return
        }
    }
}()
```
