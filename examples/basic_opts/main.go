package main

import (
	"log"
	"sync"
	"time"

	"github.com/jacobbrewer1/workerpool"
)

type task struct {
	id int

	op func() error
}

func (t *task) Run() {
	if err := t.op(); err != nil {
		// Handle error
		log.Printf("task %d failed: %v", t.id, err)
	}
}

func main() {
	// Create a new worker pool
	wp := workerpool.NewWorkerPool(
		workerpool.WithMaxQueueLength(20),
		workerpool.WithTotalWorkers(10),
		workerpool.WithDelayedStart(),
	)

	// Schedule some tasks
	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		t := &task{
			id: i,
			op: func() error {
				// Do some work
				defer wg.Done()
				log.Println("started task", i)
				time.Sleep(1 * time.Second)
				log.Println("completed task", i)
				return nil
			},
		}

		if err := wp.BlockingSchedule(t); err != nil {
			// Handle error
			log.Printf("failed to schedule task %d: %v", i, err)
		}
	}

	log.Println("Deployed workers:", wp.TotalWorkers())

	// Wait for all tasks to complete
	wg.Wait()

	log.Println("All tasks completed")

	// Stop the worker pool
	wp.Stop()
}
