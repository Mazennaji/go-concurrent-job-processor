<div align="center">

# 🚀 Go Concurrent Job Processor

**A high-performance REST API for asynchronous job processing, powered by goroutines, channels, and the worker pool pattern.**

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Architecture](https://img.shields.io/badge/Pattern-Worker_Pool-blue?style=for-the-badge)
![API](https://img.shields.io/badge/API-REST-orange?style=for-the-badge)
![No Dependencies](https://img.shields.io/badge/Dependencies-Stdlib_Only-lightgrey?style=for-the-badge)

</div>

---

## 📌 About

This project implements a **production-style concurrent job processing system** using only Go's standard library. Clients submit jobs through a REST API, jobs flow into a buffered channel queue, and a pool of goroutine workers processes them asynchronously — mirroring how real-world platforms handle background tasks like data pipelines, report generation, image processing, and email delivery.

---

## 🏗️ Architecture

```
                          ┌─────────────────────────────────────────────┐
                          │              Go HTTP Server                 │
                          │                                             │
  POST /jobs  ──────────▶ │  ┌───────────┐    ┌──────────────────────┐  │
                          │  │  Handler   │──▶ │   Job Queue (chan)   │ │
  GET /jobs/{id} ───────▶ │  └───────────┘    └──────────┬───────────┘  │
                          │                              │              │
                          │                    ┌─────────▼─────────┐    │
                          │                    │    Worker Pool    │    │
                          │                    │  ┌──┐ ┌──┐ ┌──┐   │    │
                          │                    │  │W1│ │W2│ │W3│   │    │
                          │                    │  └──┘ └──┘ └──┘   │    │
                          │                    └─────────┬─────────┘    │
                          │                              │              │
                          │                    ┌─────────▼─────────┐    │
                          │                    │   Results Store   │    │
                          │                    │   (sync.Map)      │    │
                          │                    └───────────────────┘    │
                          └─────────────────────────────────────────────┘
```

---

## 📂 Project Structure

```
go-concurrent-job-processor/
│
├── main.go          # Entry point — server startup & worker pool init
├── handlers.go      # HTTP handlers for job submission & status queries
├── worker.go        # Worker pool — goroutine lifecycle & job execution
├── job.go           # Job model, statuses, and processing logic
├── queue.go         # Buffered channel queue management
└── utils.go         # JSON response helpers & ID generation
```

---

## 📡 API Reference

### Submit a Job

```
POST /jobs
Content-Type: application/json
```

**Request:**

```json
{
  "type": "process_data"
}
```

**Response** `201 Created`:

```json
{
  "id": "a1b2c3d4",
  "type": "process_data",
  "status": "queued",
  "created_at": "2025-03-22T10:30:00Z"
}
```

---

### Check Job Status

```
GET /jobs/{id}
```

**Response** `200 OK`:

```json
{
  "id": "a1b2c3d4",
  "type": "process_data",
  "status": "completed",
  "created_at": "2025-03-22T10:30:00Z",
  "completed_at": "2025-03-22T10:30:03Z"
}
```

**Job Status Lifecycle:**

```
queued ──▶ processing ──▶ completed
                    └───▶ failed
```

---

## ⚙️ How It Works

```
1. Client sends POST /jobs           ─── Job created with status "queued"
                                          │
2. Handler pushes job into channel    ─── Buffered channel acts as the queue
                                          │
3. Worker picks up job from channel   ─── Goroutines competing for work via channel recv
                                          │
4. Worker executes job logic          ─── Status updated to "processing" → "completed"
                                          │
5. Client polls GET /jobs/{id}        ─── Retrieves final status from results store
```

---

## 🧠 Key Concepts Demonstrated

| Concept | Implementation |
|---|---|
| **Goroutines** | Each worker runs in its own goroutine, processing jobs concurrently |
| **Channels** | Buffered channel serves as a thread-safe job queue between the API layer and workers |
| **Worker Pool** | Fixed number of workers prevents resource exhaustion under high load |
| **sync.Map** | Concurrent-safe storage for job results without explicit locking |
| **Graceful Design** | Clean separation between HTTP layer, queue, and processing logic |
| **REST API** | Standard JSON API using only `net/http` from the Go stdlib |

---

## ▶️ Quick Start

**Prerequisites:** Go 1.21+

```bash
# Clone
git clone https://github.com/Mazen-Naji/go-concurrent-job-processor.git
cd go-concurrent-job-processor

# Run
go run .
```

Server starts at `http://localhost:8080`

**Test it:**

```bash
# Submit a job
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"type": "process_data"}'

# Check status (replace with actual ID from response)
curl http://localhost:8080/jobs/a1b2c3d4
```

---

## 🔌 Extensibility

Adding a new job type requires no changes to the worker pool or API layer:

```go
// Define new processing logic in job.go
func (j *Job) Execute() error {
    switch j.Type {
    case "process_data":
        return processData(j)
    case "send_email":          // ← new job type
        return sendEmail(j)
    case "generate_report":     // ← another new type
        return generateReport(j)
    default:
        return fmt.Errorf("unknown job type: %s", j.Type)
    }
}
```

---

## 🗺️ Roadmap

- [ ] Persistent storage with PostgreSQL / Redis
- [ ] Job priority queue
- [ ] Rate limiting & throttling
- [ ] API key authentication
- [ ] Docker containerization
- [ ] Real-time web dashboard for monitoring
- [ ] Distributed workers across multiple nodes
- [ ] Graceful shutdown with `context.Context`

---

## 🛠️ Tech Stack

| | |
|---|---|
| **Language** | Go 1.21+ |
| **HTTP** | `net/http` (stdlib) |
| **Concurrency** | Goroutines + Channels |
| **Pattern** | Worker Pool |
| **Storage** | `sync.Map` (in-memory) |
| **Dependencies** | None — standard library only |

---

## 🤝 Contributing

Contributions are welcome! Feel free to fork and submit a pull request.

1. Fork the repository
2. Create your feature branch — `git checkout -b feature/job-priorities`
3. Commit your changes — `git commit -m "Add job priority queue"`
4. Push to the branch — `git push origin feature/job-priorities`
5. Open a Pull Request

---

<div align="center">

**Built by [Mazen Naji](https://github.com/Mazennaji)**


⭐ If this project was useful, consider starring the repo!

</div>
