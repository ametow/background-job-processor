package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ametow/background-job-processor/internal/config"
	"github.com/ametow/background-job-processor/internal/server/domain/entity"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage() *Storage {
	db, err := sql.Open("pgx", config.DSN)

	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10)
	return &Storage{
		DB: db,
	}
}

func (s *Storage) GetNewTasks() ([]entity.TaskEntity, error) {
	log.Println("Storage.getNewTasks - hello")

	if s.DB == nil {
		return nil, errors.New("you haven`t opened the database connection")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`UPDATE tasks 
		SET task_status = $1 WHERE task_status = $2 
		RETURNING task_id, task_request_method, task_headers, task_url;`,
	)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query("in_process", "new")
	if err != nil {
		return nil, err
	}
	var newTasks []entity.TaskEntity

	for rows.Next() {
		task := entity.TaskEntity{}
		var responseHeadersBuffer []byte
		err := rows.Scan(&task.ID, &task.Method, &responseHeadersBuffer, &task.URL)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(responseHeadersBuffer, &task.Headers)
		if err != nil {
			log.Println(err)
		}

		newTasks = append(newTasks, task)
	}

	log.Println("Storage.getNewTasks got new tasks:", newTasks)
	log.Println("Storage.getNewTasks - bye")

	err = tx.Commit()

	return newTasks, err
}
func (s *Storage) BulkAddTaskResults(results []entity.ResultEntity) error {
	log.Println("Storage.BulkAddTaskResults - hello")

	if s.DB == nil {
		return errors.New("you haven`t opened the database connection")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`UPDATE tasks 
		SET task_status = $2, result_http_status_code = $3, result_headers = $4, result_body_length = $5  
    	WHERE task_id = $1;`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, r := range results {
		if _, err = stmt.Exec(r.TaskID, r.TaskStatus, r.ResponseHttpStatusCode, r.ResponseHeaders, r.ResponseBodyLength); err != nil {
			return err
		}
	}

	log.Println("Storage.BulkAddTaskResults - bye")

	return tx.Commit()
}

func (s *Storage) BulkAddTaskResultsViaCh(resultCh chan entity.ResultEntity) error {
	log.Println("Storage.BulkAddTaskResultsViaCh - hello")

	if s.DB == nil {
		return errors.New("you haven`t opened the database connection")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`UPDATE tasks 
		SET task_status = $2, result_http_status_code = $3, result_headers = $4, result_body_length = $5  
    	WHERE task_id = $1;`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

loop:
	for {
		log.Println("Storage.flushing resultCh - len is ", len(resultCh))

		select {
		case r := <-resultCh:
			if _, err = stmt.Exec(r.TaskID, r.TaskStatus, r.ResponseHttpStatusCode, r.ResponseHeaders, r.ResponseBodyLength); err != nil {
				return err
			}
		default:
			break loop
		}
	}

	log.Println("Storage.BulkAddTaskResultsViaCh - bye")

	return tx.Commit()
}
