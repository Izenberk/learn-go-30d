# Day 1 — Interfaces: implicit, small, composable

## 🧠 Why interfaces *first*?

Most languages make you declare that a type *implements* an interface (`implements Runnable`, `extends AbstractThing`). Go flips this entirely.

> **In Go, interfaces are satisfied implicitly.** If your type has the methods — it already implements the interface. No declaration needed.

This is not just syntax sugar. It's a **design philosophy**: the interface belongs to the *consumer*, not the *producer*. The package that *uses* a value defines what it needs — not the package that *provides* it.

---

## 🔬 Concept 1 — Implicit satisfaction

```go
// The interface — defined by whoever NEEDS this behavior
type Greeter interface {
    Greet() string
}

// The concrete type — knows nothing about Greeter
type EnglishSpeaker struct {
    Name string
}

func (e EnglishSpeaker) Greet() string {
    return "Hello, I'm " + e.Name
}

// This works — EnglishSpeaker satisfies Greeter implicitly
func PrintGreeting(g Greeter) {
    fmt.Println(g.Greet())
}
```

`EnglishSpeaker` never said `implements Greeter`. It just *has* the method. That's enough.

**Why this matters architecturally:** You can define an interface in your package against a type that lives in a *third-party* library — without modifying that library at all. Huge for decoupling.

---

## 🔬 Concept 2 — Keep interfaces small

The Go standard library's most-used interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)  // just 1 method
}

type Writer interface {
    Write(p []byte) (n int, err error) // just 1 method
}

type Stringer interface {
    String() string                    // just 1 method
}
```

**The rule of thumb from the Go team:**
> *"The bigger the interface, the weaker the abstraction."* — Rob Pike

One method = one capability = one reason to depend on it. This maps directly to the **Interface Segregation Principle** (the I in SOLID) — but Go enforces it culturally, not with a compiler rule.

---

## 🔬 Concept 3 — Composition over fat interfaces

Instead of one big interface, compose small ones:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose when you need both
type ReadWriter interface {
    Reader
    Writer
}
```

`io.ReadWriter` is literally defined this way in the stdlib. Your type satisfies `ReadWriter` if it has *both* methods — still implicitly, still no declaration.

---

## 💻 Code — Line-by-line walkthrough of `main.go`

### The interfaces

```go
type Storer interface {
    Store(key, value string)
}

type Fetcher interface {
    Fetch(key string) (string, bool)
}
```

Each interface declares **one capability**. `Storer` can store. `Fetcher` can fetch. That's it. Any type that has a matching method signature automatically satisfies the interface — no `implements` keyword.

```go
type Cache interface {
    Storer
    Fetcher
}
```

**Interface embedding** — `Cache` is composed from `Storer` + `Fetcher`. A type satisfies `Cache` only if it has *both* `Store()` and `Fetch()`. You only use this composed interface when a function genuinely needs both capabilities. If you don't need it, don't reach for it.

---

### The concrete type

```go
type MemCache struct {
    data map[string]string
}
```

A struct with a single field: an unexported `map`. The lowercase `data` means it's **unexported** — outside packages can't access `m.data` directly. They must go through the methods. This is Go's version of encapsulation (no `private` keyword, just casing).

---

### The constructor — and the nil map trap

```go
func NewMemCache() *MemCache {
    return &MemCache{data: make(map[string]string)}
}
```

**Why `make()` is critical here:** In Go, a `map`'s zero value is `nil`. A nil map can be *read* (returns zero values), but **writing to a nil map panics at runtime**:

```go
var m map[string]string  // m is nil
_ = m["key"]            // fine — returns ""
m["key"] = "value"      // PANIC: assignment to entry in nil map
```

`make(map[string]string)` allocates and initializes the map so it's ready for writes. This is why Go uses the **constructor pattern** (`NewXxx()`) — to guarantee the struct is in a valid state.

**Why return `*MemCache` (pointer), not `MemCache` (value)?** Because maps are reference types under the hood, but more importantly — returning a pointer means all methods with pointer receivers work without the caller needing to think about it. It also avoids copying the struct on return (relevant when structs grow larger).

---

### Pointer receivers — why `*MemCache`, not `MemCache`

```go
func (m *MemCache) Store(key, value string) {
    m.data[key] = value
}

func (m *MemCache) Fetch(key string) (string, bool) {
    v, ok := m.data[key]
    return v, ok
}
```

The `(m *MemCache)` before the function name is a **pointer receiver**. This means the method operates on the original struct, not a copy.

| Receiver | What happens | Use when |
|---|---|---|
| `(m MemCache)` — value | Method gets a **copy** of the struct | Read-only, small structs |
| `(m *MemCache)` — pointer | Method gets the **original** struct | Mutating state, or struct is large |

`Store` mutates `m.data`, so it *must* be a pointer receiver. `Fetch` only reads, but **Go convention says: if any method needs a pointer receiver, use pointer receivers for all methods on that type.** This keeps the method set consistent and avoids subtle interface satisfaction bugs.

**The subtle bug if you mix receivers:**

```go
func (m MemCache) Store(key, value string) { ... }  // value receiver
func (m *MemCache) Fetch(key string) (string, bool) { ... } // pointer receiver

var s Storer = MemCache{}   // compiles — value type has Store()
var f Fetcher = MemCache{}  // ❌ COMPILE ERROR — value type does NOT have pointer-receiver Fetch()
var f Fetcher = &MemCache{} // ✅ pointer type has both
```

A **pointer type** gets both value and pointer receiver methods. A **value type** only gets value receiver methods. Mixing leads to confusion. Don't do it.

---

### Narrow function signatures — ISP in action

```go
func Populate(s Storer) {
    s.Store("lang", "Go")
    s.Store("day", "1")
}
```

`Populate` only writes. It asks for `Storer`, not `Cache`, not `*MemCache`. This means:
- You can pass **any type** that has `Store()` — a Redis client, a mock for testing, a file-backed store
- The function's contract is explicit: "I will store things. That's all."

```go
func Display(f Fetcher) {
    if v, ok := f.Fetch("lang"); ok {
        fmt.Println("lang:", v)
    }
}
```

`Display` only reads. The **comma-ok idiom** (`v, ok := f.Fetch(...)`) is Go's pattern for "this might not exist." The second return value `ok` is `true` if the key was found, `false` otherwise. You'll see this everywhere: map lookups, type assertions, channel receives.

---

### The main function — execution flow

```go
func main() {
    c := NewMemCache()  // 1. Create a *MemCache with an initialized map
    Populate(c)         // 2. *MemCache satisfies Storer → stores "lang"="Go", "day"="1"
    Display(c)          // 3. *MemCache satisfies Fetcher → prints "lang: Go"
}
```

`c` is of type `*MemCache`. It's never declared as `Storer` or `Fetcher` — Go figures that out at the call site. This is **implicit satisfaction** in action.

**Trace the full flow:**
1. `NewMemCache()` → returns `&MemCache{data: map[string]string{}}`
2. `Populate(c)` → Go checks: does `*MemCache` have `Store(string, string)`? Yes → compiles. Stores two key-value pairs.
3. `Display(c)` → Go checks: does `*MemCache` have `Fetch(string) (string, bool)`? Yes → compiles. Fetches `"lang"`, finds it, prints `"lang: Go"`.

---

## 🔍 Quick reference — key patterns from Day 1

| Pattern | What it means | Why it matters |
|---|---|---|
| Implicit interface satisfaction | No `implements` keyword | Decouples producer from consumer — you can interface over third-party types |
| Small interfaces (1-2 methods) | `Storer`, `Fetcher` | Easier to satisfy, mock, and compose — stronger abstraction |
| Interface composition | `Cache` embeds `Storer` + `Fetcher` | Build bigger contracts from small pieces, not the other way around |
| Pointer receivers everywhere | `(m *MemCache)` | Consistent method set, avoids interface satisfaction bugs |
| `make()` for maps | `make(map[string]string)` | Nil maps panic on write — always initialize |
| Constructor pattern | `NewMemCache() *MemCache` | Returns concrete type, guarantees valid state |
| Comma-ok idiom | `v, ok := m[key]` | Safe "exists?" check — used across maps, type assertions, channels |
| Accept interfaces, return concretes | `Populate(s Storer)` / `NewMemCache() *MemCache` | Keeps callers decoupled, constructors flexible |

---

## ⚡ The one anti-pattern to burn into memory

```go
// ❌ Don't do this — accept interfaces, return concretes
func NewMemCache() Cache { ... }      // too early abstraction

// ✅ Do this
func NewMemCache() *MemCache { ... }  // concrete return
func Populate(s Storer) { ... }       // interface parameter
```

> **"Accept interfaces, return concrete types."**
> This is idiomatic Go. It keeps constructors flexible and callers decoupled.

---

## ✅ Day 1 Checkpoint

Before we move on, can you answer these in your own words?

1. Why doesn't Go need an `implements` keyword?
2. What's the risk of a 10-method interface vs a 1-method interface?
3. Why does `Populate` take `Storer` instead of `*MemCache`?
