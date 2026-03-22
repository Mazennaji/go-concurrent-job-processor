package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CreateJobRequest struct {
	Type string `json:"type"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandlers(queue *JobQueue) *Handlers {
	return &Handlers{queue: queue}
}

type Handlers struct {
	queue *JobQueue
}

func (h *Handlers) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	if req.Type == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "job type is required"})
		return
	}

	validTypes := map[string]bool{
		"process_data":    true,
		"send_email":      true,
		"generate_report": true,
	}
	if !validTypes[req.Type] {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid job type. Supported: process_data, send_email, generate_report",
		})
		return
	}

	job := NewJob(req.Type)

	if err := h.queue.Enqueue(job); err != nil {
		WriteJSON(w, http.StatusServiceUnavailable, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, job)
}

func (h *Handlers) GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "job ID is required"})
		return
	}
	id := parts[1]

	job, found := h.queue.GetJob(id)
	if !found {
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "job not found"})
		return
	}

	WriteJSON(w, http.StatusOK, job)
}

func (h *Handlers) ListJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	jobs := h.queue.GetAllJobs()
	if jobs == nil {
		jobs = []*Job{}
	}

	WriteJSON(w, http.StatusOK, jobs)
}

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "go-concurrent-job-processor",
	})
}

func (h *Handlers) JobsRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/jobs")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodPost:
			h.CreateJob(w, r)
		case http.MethodGet:
			h.ListJobs(w, r)
		default:
			WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		}
		return
	}

	if r.Method == http.MethodGet {
		h.GetJob(w, r)
		return
	}

	WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
}
