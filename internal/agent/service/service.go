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
	storage Storage
	request Request
}

func NewService(s Storage, r Request) Service {
	return Service{
		storage: s,
		request: r,
	}
}

// TODO move away with these global variables
var taskCh = make(chan entity.TaskEntity)
var resultCh = make(chan entity.ResultEntity, 5)

func (s Service) StartGettingNewTasks(timeInterval time.Duration) error {
	log.Println("service.StartGettingNewTasks - hello")
	startTimer := time.Now()

	s.createWorkers(1000)

	for i := 0; i < 10; i++ {

		err := s.getNewTasks()
		if err != nil {
			return err
		}

		time.Sleep(1 * timeInterval)
		log.Println("service.StartGettingNewTasks - resultCh - len is ", len(resultCh))

		if len(resultCh) > 0 {
			s.flushToDB()
		}

		duration := time.Since(startTimer)
		fmt.Printf("servis.getNewTasks - Execution Time ms %d\n", duration.Milliseconds())
	}
	log.Println("service.StartGettingNewTasks - bye")
	return nil
}

func (s Service) createWorkers(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for task := range taskCh {
				s.makeRequest(task)
			}
		}()
	}
}

func (s Service) makeRequest(task entity.TaskEntity) {
	log.Println("service.makeRequest - Hello")

	result, err := s.request.Request(task)
	if err != nil {
		log.Println(err)
	}
	select {
	case resultCh <- result:
		log.Println("service.makeRequest - taskCh len is ", len(taskCh))
	default:
		s.flushToDB()
		resultCh <- result
	}

	log.Println("service.makeRequest - Bye")
}

func (s Service) flushToDB() {
	log.Println("service.flushToDB - flushing resultCh - len is ", len(resultCh))
	err := s.storage.BulkAddTaskResultsViaCh(resultCh)
	if err != nil {
		log.Println(err)
	}
}

func (s Service) getNewTasks() error {
	log.Println("service.getNewTasks - hello")
	newTasks, err := s.storage.GetNewTasks()

	if err != nil {
		return err
	}

	if len(newTasks) != 0 {
		for _, task := range newTasks {
			taskCh <- task
		}
	}

	log.Println("service.getNewTasks - bye")
	return nil
}
