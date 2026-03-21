# Day 2 — Errors: wrap, sentinel, Is/As

## 🧠 Why Go errors are different

In most languages, errors are **exceptions** — they unwind the call stack automatically, you `catch` them somewhere up the chain.

Go makes a deliberate choice: **errors are just values**. They're returned explicitly, handled explicitly, and flow through your code like any other return value.

```go
result, err := doSomething()
if err != nil {
    // handle it HERE, right now
}
```

This isn't boilerplate. It's **intentional design** — errors are part of your function's contract, visible in the signature, impossible to accidentally ignore if you care about correctness.

---

## 🔬 Concept 1 — The `error` interface

Here's the entire `error` type in Go:

```go
type error interface {
    Error() string
}
```

That's it. One method. Any type that has `Error() string` is an error. This is Day 1 interfaces in action — you're already using what you learned.

---

## 🔬 Concept 2 — Sentinel errors

A **sentinel error** is a package-level error value used as a signal.

```go
// Declaration — in the package that owns it
var ErrNotFound = errors.New("not found")

// Usage — caller checks identity
if err == ErrNotFound {
    // handle the "not found" case specifically
}
```

Think of sentinels like **named constants for failure states**. `io.EOF` is the most famous one — it signals "no more data", not a real failure.

**The rule:** sentinels are for errors callers are *expected to handle differently*. Don't create one for every error — only when the caller needs to branch on it.

---

## 🔬 Concept 3 — Wrapping errors

Raw errors lose context as they bubble up. Wrapping adds a layer of *where* and *why* at each level.

```go
// %w is the wrap verb — this is the key
return fmt.Errorf("fetchUser: %w", err)
```

The result is a chain:

```
fetchUser: db query: connection refused
```

Each layer adds context without destroying the original error. The caller can still inspect the root cause — that's what `Is` and `As` are for.

**Why `%w` and not `%v`?**

```go
// %v — formats the error as a string. Chain is BROKEN.
return fmt.Errorf("fetchUser: %v", err)  // original err is lost

// %w — wraps the error. Chain is PRESERVED.
return fmt.Errorf("fetchUser: %w", err)  // original err still inspectable
```

`%v` is a one-way street. `%w` keeps the door open for `Is`/`As`.

---

## 🔬 Concept 4 — `errors.Is` and `errors.As`

### `errors.Is` — checks identity through the chain

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("fetchUser: %w", ErrNotFound)

// Direct == comparison FAILS — err is a wrapped error, not ErrNotFound itself
err == ErrNotFound  // false ❌

// errors.Is unwraps the chain until it finds ErrNotFound
errors.Is(err, ErrNotFound)  // true ✅
```

Use `errors.Is` whenever you want to check *if a specific sentinel is anywhere in the error chain*.

---

### `errors.As` — extracts a concrete type through the chain

```go
// A custom error type with extra fields
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Somewhere deep in the call stack
func validate(input string) error {
    return &ValidationError{Field: "email", Message: "invalid format"}
}

// Caller wraps it
func process(input string) error {
    if err := validate(input); err != nil {
        return fmt.Errorf("process: %w", err)
    }
    return nil
}

// Top-level handler
func main() {
    err := process("bad-input")

    var valErr *ValidationError
    if errors.As(err, &valErr) {
        // valErr is now populated — you have the full struct
        fmt.Println("field:", valErr.Field)
        fmt.Println("msg:", valErr.Message)
    }
}
```

`errors.As` unwraps until it finds a value assignable to the target type. You get the *concrete error*, not just a string.

---

## 💻 Full example — putting it all together

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel
var ErrNotFound = errors.New("not found")

// Custom error type
type DBError struct {
    Code    int
    Message string
}

func (e *DBError) Error() string {
    return fmt.Sprintf("db error %d: %s", e.Code, e.Message)
}

// Simulated DB layer
func queryDB(id int) error {
    if id == 0 {
        return ErrNotFound
    }
    if id < 0 {
        return &DBError{Code: 500, Message: "connection refused"}
    }
    return nil
}

// Service layer — wraps the error with context
func getUser(id int) error {
    if err := queryDB(id); err != nil {
        return fmt.Errorf("getUser(id=%d): %w", id, err)
    }
    return nil
}

func main() {
    // Case 1 — sentinel check
    err := getUser(0)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("handled: user not found")
    }

    // Case 2 — type extraction
    err = getUser(-1)
    var dbErr *DBError
    if errors.As(err, &dbErr) {
        fmt.Printf("handled: db code=%d msg=%s\n", dbErr.Code, dbErr.Message)
    }

    // Case 3 — full error string shows the chain
    fmt.Println(err)
    // → getUser(id=-1): db error 500: connection refused
}
```

---

## 🗺️ Mental model — when to use what

| Tool | Use when |
|---|---|
| `errors.New` | Simple sentinel — caller branches on identity |
| `fmt.Errorf("%w")` | Adding context at each layer |
| Custom type | Caller needs structured data from the error |
| `errors.Is` | Checking if a sentinel is in the chain |
| `errors.As` | Extracting a typed error from the chain |

---

## ⚡ The anti-patterns to avoid

```go
// ❌ Wrapping with %v — breaks the chain
return fmt.Errorf("something failed: %v", err)

// ❌ Bare string comparison — fragile
if err.Error() == "not found" { ... }

// ❌ Swallowing errors — silent failure
result, _ := doSomething()

// ✅ Wrap with %w, check with Is/As
return fmt.Errorf("getUser: %w", err)
errors.Is(err, ErrNotFound)
```

---

## 🔗 Aegis Stream connection

You're already doing this in production. In Aegis Stream, when a TCP connection drops or a NATS publish fails — those errors bubble up through your worker pool. Wrapping at each layer (`tcpHandler: %w`, `worker: %w`) gives you a full trace in logs without losing the root cause. `errors.Is` lets your retry logic say *"is this a transient network error or a fatal config error?"* — and branch accordingly.

---

## ✅ Day 2 Checkpoint

Answer in your own words:

1. What's the difference between `%w` and `%v` in `fmt.Errorf`?
2. When would you use `errors.Is` vs `errors.As`?
3. Why are sentinel errors package-level variables and not just strings?
