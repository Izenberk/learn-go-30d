# Day 01 — Interface Syntax

## Defining an interface
```go
type Greeter interface {
    Greet() string    // method name + signature — no body, no implementation
}
```
Plain English: "Any type that has a `Greet()` method returning a string satisfies this interface."

## Implementing an interface (implicit)
```go
type EnglishSpeaker struct{ Name string }

func (e EnglishSpeaker) Greet() string {
    return "Hello, I'm " + e.Name
}
```
Plain English: No `implements` keyword needed. `EnglishSpeaker` has `Greet()` — it automatically satisfies `Greeter`.

## Using an interface as a function parameter
```go
func PrintGreeting(g Greeter) {
    fmt.Println(g.Greet())
}
```
Plain English: `PrintGreeting` accepts any value that satisfies `Greeter` — it doesn't care about the concrete type.

## Composing interfaces
```go
type Reader interface { Read(p []byte) (int, error) }
type Writer interface { Write(p []byte) (int, error) }

type ReadWriter interface {
    Reader   // embed Reader
    Writer   // embed Writer
}
```
Plain English: `ReadWriter` requires both `Read` and `Write`. A type must have both methods to satisfy it.

## Constructor pattern
```go
func NewMemCache() *MemCache { ... }  // returns concrete type ✅
func NewMemCache() Cache { ... }      // returns interface ❌ — too early
```
Rule: **Accept interfaces, return concrete types.**

## Method receiver syntax
```go
func (m *MemCache) Store(key, value string) { }
//    ^^^^^^^^^
//    receiver — like "this" or "self" in other languages
//    * means pointer receiver — modifies the original struct
```
