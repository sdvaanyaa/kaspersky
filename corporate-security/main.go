package main

import (
	"flag"
	"fmt"
	"github.com/sdvaanyaa/kaspersky/corporate-security/workerpool"
	"math/rand"
	"time"
)

func main() {
	workers := flag.Int("workers", 4, "number of worker goroutines")
	numTasks := flag.Int("tasks", 100, "number of tasks to submit")
	flag.Parse()

	wp := workerpool.New(*workers)

	for i := 1; i <= *numTasks; i++ {
		wp.Submit(func() {
			fmt.Printf("Task %d starting\n", i)
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // simulate work
			fmt.Printf("Task %d done\n", i)
		})
	}

	wp.StopWait() // wait for all tasks to finish (including queued ones)
	// wp.Stop()     // stop immediately, only running tasks will finish

	fmt.Println("Pool stopped")
}
