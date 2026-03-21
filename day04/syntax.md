# Day 04 — Channel Syntax

## Creating channels
```go
ch := make(chan int)       // unbuffered — capacity 0, rendezvous required
ch := make(chan int, 3)    // buffered — capacity 3, holds up to 3 values
```

## The arrow rule — direction of data flow
```go
ch <- 42        // SEND — arrow points INTO ch — "put 42 into ch"
val := <-ch     // RECEIVE — arrow points OUT of ch — "take value out of ch"
<-ch            // RECEIVE and discard — drain one value, ignore it
```

## Receive with ok check
```go
val, ok := <-ch    // ok = true  → channel open, val is valid
                   // ok = false → channel closed and empty, val is zero value
```

## Close a channel
```go
close(ch)    // signal: no more values will be sent
             // only the SENDER closes — never the receiver
             // sending to a closed channel panics
```

## Range over a channel
```go
for val := range ch {    // receives values until ch is closed
    fmt.Println(val)     // no ok check needed — range handles it
}
```
Plain English: `for range` on a channel is like `for range` on a slice — it stops automatically when the source is exhausted (closed).

## Directional channel types
```go
chan<- int    // send-only — "this channel only accepts int going in"
<-chan int    // receive-only — "this channel only delivers int going out"
chan int      // bidirectional — can both send and receive
```
Arrow direction mirrors the data flow direction — same rule as send/receive.

## Nil channel behaviour
```go
var ch chan int    // zero value of channel type = nil
ch <- 1           // blocks forever — nil send never completes
<-ch              // blocks forever — nil receive never completes
```
Use case: in a select, a nil channel case is permanently disabled — useful for dynamic select control (Day 05).

## Channel in function signatures
```go
func produce(out chan<- int) { out <- 42 }    // can only send
func consume(in <-chan int) { v := <-in }     // can only receive
// bidirectional chan int auto-converts to directional at call site
```
