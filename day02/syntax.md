# Day 02 — Error Syntax

## The error interface
```go
type error interface {
    Error() string    // any type with this method IS an error
}
```

## Creating a simple error
```go
err := errors.New("something went wrong")
// returns a value that satisfies the error interface
```

## Returning an error from a function
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")  // return zero value + error
    }
    return a / b, nil  // nil = no error
}
```

## Checking an error
```go
result, err := divide(10, 0)
if err != nil {         // nil means no error
    fmt.Println(err)    // calls err.Error() internally
    return
}
```

## Sentinel error — package-level error variable
```go
var ErrNotFound = errors.New("not found")   // declared once at package level
                                            // var, not const — errors are values
```

## Wrapping an error with context
```go
return fmt.Errorf("getUser: %w", err)
//                           ^
//                           %w = wrap verb — preserves the original error
//                           %v = format verb — destroys the chain
```
Plain English: `%w` keeps a door open to inspect the original error. `%v` prints it as a string and closes the door.

## Checking wrapped errors
```go
errors.Is(err, ErrNotFound)     // is ErrNotFound anywhere in the chain?
                                // unwraps automatically — works through %w

var dbErr *DBError
errors.As(err, &dbErr)          // is there a *DBError anywhere in the chain?
                                // if yes, dbErr is populated with the actual value
//              ^
//              & = pass pointer to target — errors.As needs to assign into it
```

## Custom error type
```go
type ValidationError struct {   // a struct that satisfies the error interface
    Field   string
    Message string
}

func (e *ValidationError) Error() string {   // this method makes it an error
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}
```
