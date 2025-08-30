package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/sdvaanyaa/kaspersky/sandbox-development/pool"
)

func main() {
	workers := flag.Int("workers", 2, "number of worker goroutines")
	queueSize := flag.Int("queue", 3, "maximum queue size for tasks")
	numTasks := flag.Int("tasks", 5, "number of tasks to submit")
	flag.Parse()

	hook := func() {
		fmt.Println("Task completed")
	}

	p := pool.New(*workers, *queueSize, hook)

	for i := 1; i <= *numTasks; i++ {
		err := p.Submit(func() {
			fmt.Printf("Task %d starting\n", i)
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // simulate work
			fmt.Printf("Task %d done\n", i)
		})

		if err != nil {
			fmt.Printf("Submit error for task %d: %v\n", i, err)
		}
	}

	err := p.Stop()
	if err != nil {
		fmt.Printf("Stop error: %v\n", err)
	}

	fmt.Println("Pool stopped")
}
