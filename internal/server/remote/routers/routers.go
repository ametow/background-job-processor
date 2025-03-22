package routers

import (
	"github.com/ametow/background-job-processor/internal/server/remote/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(task handlers.TaskHandler) chi.Router {

	r := chi.NewRouter()

	r.Post("/task", task.CreateTask)
	r.Get("/task/{taskID}", task.GetTaskStatus)
	return r
}
