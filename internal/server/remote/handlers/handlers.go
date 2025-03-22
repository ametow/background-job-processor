package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ametow/background-job-processor/internal/server/domain/entity"
	localError "github.com/ametow/background-job-processor/internal/server/domain/error"
)

type TaskService interface {
	CreateTask(te entity.TaskEntity)
	GetTaskStatus(id string) (entity.ResultEntity, error)
}

type TaskHandler struct {
	taskService TaskService
}

func NewTaskHandler(ts TaskService) TaskHandler {
	return TaskHandler{
		taskService: ts,
	}
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

const IDLength int = 4

func (th TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler CreateTask - hello")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, localError.WrongContentType.Error(), http.StatusBadRequest)
		return
	}

	var decVal struct { // decoded value
		Method  string
		URL     string
		Headers entity.Headers //http.Header
	}

	if err := json.NewDecoder(r.Body).Decode(&decVal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Handler CreateTask - request: Task", decVal)

	bytesTaskID, err := generateRandom(IDLength)
	if err != nil {
		log.Println(err)
		return
	}
	taskID := hex.EncodeToString(bytesTaskID)

	te := entity.TaskEntity{
		ID:      taskID,
		Method:  decVal.Method,
		URL:     decVal.URL,
		Headers: decVal.Headers,
	}
	th.taskService.CreateTask(te)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	encVal := struct { // encoded value
		ID string `json:"id"`
	}{
		ID: taskID,
	}

	log.Println("Handler CreateTask - bye")
	json.NewEncoder(w).Encode(encVal)

}

func (th TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler GetTaskStatus - hello")
	taskID := chi.URLParam(r, "taskID")

	log.Println("Handler GetTaskStatus, taskID", taskID)

	re, err := th.taskService.GetTaskStatus(taskID)

	if err == localError.NotFound {
		http.Error(w, localError.NotFound.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Handler GetTaskStatus - bye")
	json.NewEncoder(w).Encode(re)
}
