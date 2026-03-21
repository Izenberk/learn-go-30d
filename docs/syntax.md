# Go Universal Syntax Reference

## Variable declaration
```go
var x int          // declare with zero value
x := 42            // short declare + assign (only inside functions)
var x int = 42     // explicit type + value
const Pi = 3.14    // constant — immutable
```

## Functions
```go
func name(param type) returnType { }        // basic
func name(a, b int) (int, error) { }        // multiple returns
func name(a int) (result int, err error) { } // named returns
```

## Structs and methods
```go
type Dog struct {                 // define a struct
    Name string
    Age  int
}

func (d Dog) Speak() string { }   // value receiver — d is a copy
func (d *Dog) Grow() { }          // pointer receiver — d is the original
```

## Interfaces
```go
type Writer interface {           // define an interface
    Write(p []byte) (int, error)  // method signature only — no body
}
```

## Pointers
```go
x := 42
p := &x     // & = "address of x" — p holds x's memory address
*p = 99     // * = "value at address" — changes x to 99
```

## Error handling
```go
val, err := doSomething()   // functions return error as last value
if err != nil {             // always check immediately
    return err
}
```

## Control flow
```go
if x > 0 { }                           // no parentheses needed
for i := 0; i < 10; i++ { }            // standard loop
for condition { }                       // while-style loop
for { }                                 // infinite loop
for i, v := range slice { }            // range over slice
for k, v := range map { }              // range over map
switch x { case 1: ... case 2: ... }   // no fallthrough by default
```

## Types
```go
// Zero values — what you get when declared but not assigned
int       → 0
string    → ""
bool      → false
pointer   → nil
slice     → nil
map       → nil
channel   → nil
```

## Make vs New
```go
make(chan int)        // for channels, slices, maps — initializes internals
make([]int, 5)        // slice of length 5
make(map[string]int)  // empty map ready to use
new(int)              // allocates zeroed memory, returns pointer — rarely used
```

## Defer
```go
defer f()    // f() runs when the surrounding function returns
             // always — even on panic or early return
             // multiple defers execute in LIFO order (last in, first out)
```
