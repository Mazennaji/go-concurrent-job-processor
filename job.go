package main

import (
	"fmt"
	"math/rand"
	"time"
)

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Status      JobStatus  `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
}

func NewJob(jobType string) *Job {
	return &Job{
		ID:        GenerateID(),
		Type:      jobType,
		Status:    StatusQueued,
		CreatedAt: time.Now().UTC(),
	}
}

func (j *Job) Execute() error {
	switch j.Type {
	case "process_data":
		return processData(j)
	case "send_email":
		return sendEmail(j)
	case "generate_report":
		return generateReport(j)
	default:
		return fmt.Errorf("unknown job type: %s", j.Type)
	}
}

func processData(j *Job) error {
	duration := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(duration)
	fmt.Printf("[Worker] Job %s: Data processed successfully (took %v)\n", j.ID, duration)
	return nil
}

func sendEmail(j *Job) error {
	duration := time.Duration(500+rand.Intn(1500)) * time.Millisecond
	time.Sleep(duration)
	fmt.Printf("[Worker] Job %s: Email sent successfully (took %v)\n", j.ID, duration)
	return nil
}

func generateReport(j *Job) error {
	duration := time.Duration(2+rand.Intn(4)) * time.Second
	time.Sleep(duration)
	fmt.Printf("[Worker] Job %s: Report generated successfully (took %v)\n", j.ID, duration)
	return nil
}
