package main

import (
	"errors"
	"fmt"
)

// Sentinel
var ErrNotFound = errors.New("not found")

// Custom error type
type DBError struct	{
	Code 	int
	Message	string
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
}