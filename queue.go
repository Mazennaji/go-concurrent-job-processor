package main

var JobQueue chan *Job

func InitQueue(size int) {
	JobQueue = make(chan *Job, size)
}
