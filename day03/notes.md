# Day 3 — Goroutines: GMP scheduler, leaks

## 🧠 Why goroutines are a big deal

Every language has threads. Go has goroutines — and they are **not the same thing**.

A thread is an OS-managed unit. Heavy to create (~1MB stack), slow to context-switch, limited in count (thousands, maybe).

A goroutine is a **Go runtime-managed unit**. Starts at 2KB of stack, grows dynamically, context-switches in user space. You can run **hundreds of thousands** of them concurrently on a handful of OS threads.

This is why Aegis Stream can handle ~239k events/sec — the worker pool is goroutine-based, not thread-based.

---

## 🔬 Concept 1 — The GMP Scheduler

This is the engine under the hood. Three components:

```
G — Goroutine    the unit of work (your goroutine)
M — Machine      an OS thread (managed by the OS)
P — Processor    a logical CPU slot (GOMAXPROCS controls how many)
```

The relationship:

```
P1 [ G1 → G2 → G3 ]  ← run queue
     |
     M1  ← OS thread executing G1 right now

P2 [ G4 → G5 ]
     |
     M2  ← OS thread executing G4 right now
```

A `G` can only run when it's attached to a `P`. A `P` can only execute on an `M`. By default, `GOMAXPROCS = number of CPU cores` — so you get true parallelism up to your core count, with the scheduler multiplexing thousands of goroutines across that pool.

**Work stealing** — if P1's queue is empty, it steals goroutines from P2's queue. This keeps all processors busy without you writing a single line of scheduling code.

**Preemption** — since Go 1.14, goroutines are preemptible at any safe point. Long-running goroutines don't starve others.

---

## 🔬 Concept 2 — Spawning a goroutine

```go
go someFunction()       // fire and forget
go func() { ... }()     // anonymous goroutine
```

That's the entire syntax. One keyword.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("goroutine running")
    }()

    // ❌ without this, main exits before the goroutine runs
    time.Sleep(10 * time.Millisecond)
}
```

The `time.Sleep` is a smell — we'll fix it with `sync.WaitGroup` on Day 7. But for now, understand why it's needed: **`main` exiting kills all goroutines instantly**, regardless of whether they've finished.

---

## 🔬 Concept 3 — Goroutine leaks

A goroutine leak is a goroutine that is **started but never exits**. It sits in memory forever, consuming resources silently. This is one of the most common production bugs in Go.

### How leaks happen

**Pattern 1 — blocked on a channel receive forever:**

```go
func leak() {
    ch := make(chan int) // unbuffered, nobody sends
    go func() {
        val := <-ch // blocks here forever — leaked
        fmt.Println(val)
    }()
    // function returns, but goroutine is still alive
}
```

**Pattern 2 — blocked on a channel send with no receiver:**

```go
func leak() {
    ch := make(chan int) // unbuffered, nobody receives
    go func() {
        ch <- 42 // blocks here forever — leaked
    }()
}
```

**Pattern 3 — missing done signal:**

```go
func leak(ch <-chan int) {
    go func() {
        for v := range ch {
            process(v)
        }
        // only exits when ch is closed
        // if ch is never closed — goroutine lives forever
    }()
}
```

---

## 🔬 Concept 4 — The fix: always give goroutines an exit path

The pattern is **context cancellation** (Day 6 goes deep on this) combined with a select:

```go
func noLeak(ctx context.Context, ch <-chan int) {
    go func() {
        for {
            select {
            case v := <-ch:
                process(v)
            case <-ctx.Done(): // exit path — always present
                return
            }
        }
    }()
}
```

Every goroutine you spawn should have a clear answer to: **"how does this goroutine exit?"**

---

## 💻 Full example — spot the leak, then fix it

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "time"
)

// ❌ Leaky version
func startLeaky() {
    ch := make(chan int)
    go func() {
        // nobody ever sends to ch — this goroutine is stuck forever
        val := <-ch
        fmt.Println("got:", val)
    }()
}

// ✅ Fixed version — goroutine has an exit path via context
func startClean(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println("got:", val)
        case <-ctx.Done():
            fmt.Println("goroutine exiting cleanly")
            return
        }
    }()
}

func main() {
    fmt.Println("goroutines before:", runtime.NumGoroutine())

    // Spawn 5 leaky goroutines
    for i := 0; i < 5; i++ {
        startLeaky()
    }

    fmt.Println("goroutines after leaky:", runtime.NumGoroutine())
    // → 6 (main + 5 leaked)

    // Spawn 5 clean goroutines and cancel them
    ctx, cancel := context.WithCancel(context.Background())
    for i := 0; i < 5; i++ {
        startClean(ctx)
    }
    cancel() // signals all clean goroutines to exit

    time.Sleep(10 * time.Millisecond) // let them exit
    fmt.Println("goroutines after clean:", runtime.NumGoroutine())
    // → 6 still (main + 5 leaked — the clean ones exited)
}
```

`runtime.NumGoroutine()` is your leak detector in tests. If the count keeps climbing — you have a leak.

---

## 🔗 Health-Check API connection

In your Health-Check API, every concurrent health check probe will be a goroutine. The pattern you'll use:

```go
// Each target gets a goroutine, context controls the timeout
func checkAll(ctx context.Context, targets []string) []Result {
    results := make([]Result, len(targets))
    var wg sync.WaitGroup

    for i, target := range targets {
        wg.Add(1)
        go func(i int, target string) {
            defer wg.Done()
            results[i] = probe(ctx, target) // ctx carries the deadline
        }(i, target)
    }

    wg.Wait()
    return results
}
```

No context → no exit path → leak under timeout. This is exactly the bug that hits production health checkers first.

---

## 🔍 How to detect leaks in tests

```go
func TestNoLeak(t *testing.T) {
    before := runtime.NumGoroutine()

    ctx, cancel := context.WithCancel(context.Background())
    startClean(ctx)
    cancel()

    time.Sleep(10 * time.Millisecond)
    after := runtime.NumGoroutine()

    if after > before {
        t.Errorf("goroutine leak: before=%d after=%d", before, after)
    }
}
```

Or use the `goleak` package by Uber — it does this automatically at test teardown.

---

## ✅ Day 3 Checkpoint

1. What are G, M, and P in the GMP scheduler — in your own words?
2. Why does `main` exiting kill all goroutines?
3. What's the common root cause of a goroutine leak?
4. How does passing a `context.Context` to a goroutine prevent a leak?
