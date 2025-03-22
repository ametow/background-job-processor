package main

import (
	"log"
	"time"

	"github.com/ametow/background-job-processor/internal/agent/remote"
	"github.com/ametow/background-job-processor/internal/agent/service"
	"github.com/ametow/background-job-processor/internal/agent/storage"
)

func main() {
	storage := storage.NewStorage()
	remoteRequest := remote.NewRequest()
	service := service.NewService(storage, remoteRequest)
	newTaskCheckInterval := 500 * time.Millisecond

	log.Println("Getting new tasks from storage...")
	service.StartGettingNewTasks(newTaskCheckInterval)
}
