package main

import "time"

type JobStatus string

const (
	Queued     JobStatus = "queued"
	Processing JobStatus = "processing"
	Completed  JobStatus = "completed"
	Failed     JobStatus = "failed"
)

type Job struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Status      JobStatus  `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (j *Job) Execute() error {
	// simulate work
	time.Sleep(3 * time.Second)
	return nil
}
