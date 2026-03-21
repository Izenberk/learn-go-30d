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

	// Spawn 5 clean goroutine and cancel them
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 5; i++ {
		startClean(ctx)
	}
	fmt.Println("after spawning clean:", runtime.NumGoroutine()) // → 11

	cancel()	// signals all clean goroutines to exit

	time.Sleep(10 * time.Millisecond) // let them exit
	fmt.Println("goroutines after clean:", runtime.NumGoroutine())
	// → 6 still (main + 5 leaked — the clean ones exited)
}