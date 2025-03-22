package service

import (
	"fmt"
	"log"
	"time"

	"github.com/ametow/background-job-processor/internal/server/domain/entity"
)

type Storage interface {
	GetNewTasks() ([]entity.TaskEntity, error)
	BulkAddTaskResultsViaCh(resultCh chan entity.ResultEntity) error
}

type Request interface {
	Request(taskEntity entity.TaskEntity) (entity.ResultEntity, error)
}

type Service struct {
	storage       Storage
	request       Request
	taskChannel   chan entity.TaskEntity
	resultChannel chan entity.ResultEntity
}

func New(s Storage, r Request) Service {
	return Service{
		storage:       s,
		request:       r,
		taskChannel:   make(chan entity.TaskEntity),
		resultChannel: make(chan entity.ResultEntity, 5),
	}
}

func (s Service) Start(timeInterval time.Duration) error {
	log.Println("service.StartGettingNewTasks - start")
	startTimer := time.Now()

	s.createWorkers(1000)

	for {
		err := s.getNewTasks()
		if err != nil {
			return err
		}

		time.Sleep(timeInterval)
		log.Println("service.StartGettingNewTasks - s.resultChannel - len is ", len(s.resultChannel))

		if len(s.resultChannel) > 0 {
			s.flushToDB()
		}

		duration := time.Since(startTimer)
		fmt.Printf("servis.getNewTasks - Execution Time ms %d\n", duration.Milliseconds())
	}
	return nil
}

func (s Service) createWorkers(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for task := range s.taskChannel {
				s.makeRequest(task)
			}
		}()
	}
}

func (s Service) makeRequest(task entity.TaskEntity) {
	log.Println("service.makeRequest - start")

	result, err := s.request.Request(task)
	if err != nil {
		log.Println(err)
	}
	select {
	case s.resultChannel <- result:
		log.Println("service.makeRequest - s.taskChannel len is ", len(s.taskChannel))
	default:
		s.flushToDB()
		s.resultChannel <- result
	}

	log.Println("service.makeRequest - end")
}

func (s Service) flushToDB() {
	log.Println("service.flushToDB - flushing s.resultChannel - len is ", len(s.resultChannel))
	err := s.storage.BulkAddTaskResultsViaCh(s.resultChannel)
	if err != nil {
		log.Println(err)
	}
}

func (s Service) getNewTasks() error {
	log.Println("service.getNewTasks - start")
	newTasks, err := s.storage.GetNewTasks()

	if err != nil {
		return err
	}

	if len(newTasks) != 0 {
		for _, task := range newTasks {
			s.taskChannel <- task
		}
	}

	log.Println("service.getNewTasks - end")
	return nil
}
