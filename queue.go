package main

import (
	"fmt"
	"sync"
)

type JobQueue struct {
	jobs  chan *Job
	store sync.Map
	size  int
}

func NewJobQueue(bufferSize int) *JobQueue {
	return &JobQueue{
		jobs: make(chan *Job, bufferSize),
		size: bufferSize,
	}
}

func (q *JobQueue) Enqueue(job *Job) error {
	select {
	case q.jobs <- job:
		q.store.Store(job.ID, job)
		fmt.Printf("[Queue] Job %s enqueued (type: %s)\n", job.ID, job.Type)
		return nil
	default:
		return fmt.Errorf("job queue is full (capacity: %d)", q.size)
	}
}

func (q *JobQueue) Dequeue() <-chan *Job {
	return q.jobs
}

func (q *JobQueue) GetJob(id string) (*Job, bool) {
	value, ok := q.store.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*Job), true
}

func (q *JobQueue) UpdateJob(job *Job) {
	q.store.Store(job.ID, job)
}

func (q *JobQueue) GetAllJobs() []*Job {
	var jobs []*Job
	q.store.Range(func(key, value interface{}) bool {
		jobs = append(jobs, value.(*Job))
		return true
	})
	return jobs
}
