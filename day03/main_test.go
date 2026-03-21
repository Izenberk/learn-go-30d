package main

import (
    "context"
    "runtime"
    "testing"
    "time"
)

func TestLeakyGoroutine(t *testing.T) {
    before := runtime.NumGoroutine()

    startLeaky()

    time.Sleep(10 * time.Millisecond)
    after := runtime.NumGoroutine()

    if after > before {
        t.Errorf("goroutine leak: before=%d after=%d", before, after)
    }
}

func TestCleanGoroutine(t *testing.T) {
    before := runtime.NumGoroutine()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    startClean(ctx)
    cancel()

    time.Sleep(10 * time.Millisecond)
    after := runtime.NumGoroutine()

    if after > before {
        t.Errorf("goroutine leak: before=%d after=%d", before, after)
    }
}