package service

import (
	"log"

	"github.com/ametow/background-job-processor/internal/server/domain/entity"
)

type TaskStorage interface {
	CreateTask(te entity.TaskEntity) error
	GetTaskStatus(taskID string) (entity.ResultEntity, error)
}

type TaskService struct {
	taskStorage TaskStorage
}

func NewTaskService(taskStorage TaskStorage) TaskService {
	return TaskService{
		taskStorage: taskStorage,
	}
}

func (taskService TaskService) CreateTask(te entity.TaskEntity) {
	log.Println("Service CreateTask - start")

	err := taskService.taskStorage.CreateTask(te)
	if err != nil {
		log.Println(err)
	}

	log.Println("Service CreateTask - end")
}

func (taskService TaskService) GetTaskStatus(taskID string) (entity.ResultEntity, error) {
	log.Println("Service GetTaskStatus - start")

	re, err := taskService.taskStorage.GetTaskStatus(taskID)
	if err != nil {
		log.Println(err)
	}

	log.Println("Service GetTaskStatus - end")
	return re, err
}
