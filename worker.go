package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	workerCount int
	queue       *JobQueue
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int, queue *JobQueue) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount: workerCount,
		queue:       queue,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (wp *WorkerPool) Start() {
	fmt.Printf("[Pool] Starting %d workers...\n", wp.workerCount)
	for i := 1; i <= wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	fmt.Printf("[Pool] All %d workers are running\n", wp.workerCount)
}

func (wp *WorkerPool) Stop() {
	fmt.Println("[Pool] Shutting down workers...")
	wp.cancel()
	wp.wg.Wait()
	fmt.Println("[Pool] All workers stopped")
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	fmt.Printf("[Worker %d] Ready and waiting for jobs\n", id)

	for {
		select {
		case <-wp.ctx.Done():
			fmt.Printf("[Worker %d] Shutting down\n", id)
			return
		case job, ok := <-wp.queue.Dequeue():
			if !ok {
				fmt.Printf("[Worker %d] Queue closed, shutting down\n", id)
				return
			}
			wp.processJob(id, job)
		}
	}
}

func (wp *WorkerPool) processJob(workerID int, job *Job) {
	fmt.Printf("[Worker %d] Processing job %s (type: %s)\n", workerID, job.ID, job.Type)

	job.Status = StatusProcessing
	wp.queue.UpdateJob(job)

	err := job.Execute()

	now := time.Now().UTC()
	job.CompletedAt = &now

	if err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		fmt.Printf("[Worker %d] Job %s failed: %v\n", workerID, job.ID, err)
	} else {
		job.Status = StatusCompleted
		fmt.Printf("[Worker %d] Job %s completed successfully\n", workerID, job.ID)
	}

	wp.queue.UpdateJob(job)
}
