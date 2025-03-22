package main

import (
	"log"
	"net/http"

	"github.com/ametow/background-job-processor/internal/server/remote/handlers"
	"github.com/ametow/background-job-processor/internal/server/remote/routers"
	"github.com/ametow/background-job-processor/internal/server/service"
	"github.com/ametow/background-job-processor/internal/server/storage"
)

func main() {
	taskStorage := storage.NewTaskStorage()
	taskService := service.NewTaskService(taskStorage)
	taskHandler := handlers.NewTaskHandler(taskService)
	r := routers.NewRouter(taskHandler)
	serverAddress := ":8080"
	log.Println("Serving on port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
