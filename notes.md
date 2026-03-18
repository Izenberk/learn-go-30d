# Learn Go in 30 Days — Quick Reference Notes

Personal notes. Written to be scanned fast and understood by a Go newbie.

---

## Day 1 — Interfaces: implicit, small, composable

### Core idea

Go interfaces are satisfied **implicitly**. No `implements` keyword. If your type has the method, it already satisfies the interface.

```go
type Greeter interface {
    Greet() string
}

type Dog struct{ Name string }

// Dog never says "implements Greeter" — it just has the method
func (d Dog) Greet() string { return "Woof, I'm " + d.Name }
```

### Why this matters

- The **consumer** defines the interface, not the producer
- You can create an interface for a third-party type without touching their code
- Decoupling happens naturally

### Keep interfaces small (1-2 methods)

```go
type Reader interface { Read(p []byte) (int, error) }   // 1 method
type Writer interface { Write(p []byte) (int, error) }   // 1 method
```

> *"The bigger the interface, the weaker the abstraction."* — Rob Pike

Compose when you need more:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

### Pointer vs value receivers

| Receiver | Effect | Use when |
|---|---|---|
| `(m MemCache)` value | gets a **copy** | read-only, small struct |
| `(m *MemCache)` pointer | gets the **original** | mutating state, or struct is large |

**Rule:** if *any* method needs a pointer receiver, use pointer receivers for *all* methods on that type. Mixing causes interface satisfaction bugs.

### Key patterns

| Pattern | Example | Why |
|---|---|---|
| Accept interfaces, return concretes | `func Save(s Storer)` / `func New() *Thing` | callers stay decoupled |
| Constructor pattern | `func NewMemCache() *MemCache` | guarantees valid state (e.g. `make(map)`) |
| Comma-ok idiom | `v, ok := m["key"]` | safe "does it exist?" check |
| `make()` for maps | `make(map[string]string)` | nil map panics on write |

---

## Day 2 — Errors: wrap, sentinel, Is/As

### Core idea

Go has **no exceptions**. Errors are just values — returned explicitly, handled explicitly.

```go
result, err := doSomething()
if err != nil {
    // handle it right here
}
```

The `error` interface is one method:

```go
type error interface {
    Error() string
}
```

Any type with `Error() string` is an error. That's it.

### Sentinel errors — named error constants

A package-level `var` used as a known signal callers can check against.

```go
var ErrNotFound = errors.New("not found")
```

Only create sentinels when callers need to **branch** on a specific failure. Don't make one for every error.

### Custom error types — errors with data

When callers need more than just a string (e.g. error codes, field names):

```go
type DBError struct {
    Code    int
    Message string
}

func (e *DBError) Error() string {
    return fmt.Sprintf("db error %d: %s", e.Code, e.Message)
}
```

### Wrapping — adding context without losing the cause

```go
return fmt.Errorf("getUser(id=%d): %w", id, err)
// output: getUser(id=42): not found
```

**`%w` vs `%v` — this is critical:**

| Verb | What happens | Chain |
|---|---|---|
| `%w` | wraps the error | **preserved** — `Is`/`As` can still find the original |
| `%v` | formats as string | **broken** — original error is lost forever |

Always use `%w` when wrapping.

### `errors.Is` — "is this sentinel anywhere in the chain?"

```go
err := getUser(0)
// err = "getUser(id=0): not found"  (wrapped)

err == ErrNotFound           // false — err is the WRAPPED version
errors.Is(err, ErrNotFound)  // true  — unwraps the chain and finds it
```

### `errors.As` — "extract the typed error from the chain"

```go
var dbErr *DBError
if errors.As(err, &dbErr) {
    // dbErr is populated with the actual struct
    fmt.Println(dbErr.Code, dbErr.Message)
}
```

`Is` checks identity (sentinel). `As` extracts data (custom type).

### When to use what

| Tool | Use when |
|---|---|
| `errors.New("msg")` | simple sentinel — caller branches on identity |
| `fmt.Errorf("context: %w", err)` | adding context at each layer |
| custom error struct | caller needs structured data (codes, fields) |
| `errors.Is(err, target)` | checking if a sentinel is in the chain |
| `errors.As(err, &target)` | extracting a typed error from the chain |

### Anti-patterns

```go
// BAD: %v breaks the chain
return fmt.Errorf("failed: %v", err)

// BAD: string comparison is fragile
if err.Error() == "not found" { ... }

// BAD: swallowing errors — silent failure
result, _ := doSomething()

// GOOD: wrap with %w, inspect with Is/As
return fmt.Errorf("getUser: %w", err)
```
