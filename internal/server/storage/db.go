package storage

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ametow/background-job-processor/internal/config"
	"github.com/ametow/background-job-processor/internal/server/domain/entity"
	localError "github.com/ametow/background-job-processor/internal/server/domain/error"
)

type TaskStorage struct {
	DB *sql.DB
}

func NewTaskStorage() *TaskStorage {
	db, err := sql.Open("pgx", config.DSN)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(5)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks(
	TaskID serial PRIMARY KEY,
	task_id text,
	task_request_method text,
	task_url text,
	task_headers json,
	task_status text DEFAULT 'new',
	result_http_status_code INTEGER,
	result_headers json,
	result_body_length text
	                         );`)
	if err != nil {
		log.Fatal(err)
	}

	return &TaskStorage{
		DB: db,
	}
}

func (ts TaskStorage) CreateTask(te entity.TaskEntity) error {
	log.Println("Storage CreateTask - hello")

	_, err := ts.DB.Exec(
		`INSERT INTO tasks (task_id , task_request_method , task_url, task_headers)
		VALUES ($1, $2, $3, $4);`,
		te.ID, te.Method, te.URL, te.Headers,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Storage CreateTask - bye")

	return nil
}

func (ts TaskStorage) GetTaskStatus(taskID string) (entity.ResultEntity, error) {
	log.Println("Storage GetTaskStatus - hello")

	row := ts.DB.QueryRow(
		`SELECT task_status, result_http_status_code, result_headers, result_body_length
		FROM tasks WHERE task_id = $1;`,
		taskID)

	re := entity.ResultEntity{}
	var responseHeadersBuffer []byte
	err := row.Scan(&re.TaskStatus, &re.ResponseHttpStatusCode, &responseHeadersBuffer, &re.ResponseBodyLength)
	if err != nil {
		log.Println(err)
	}

	if err == sql.ErrNoRows {
		log.Println("Storage GetTaskStatus, record not found")
		return entity.ResultEntity{}, localError.NotFound //errors.New("Resource was not found")
	}
	err = json.Unmarshal(responseHeadersBuffer, &re.ResponseHeaders)
	if err != nil {
		log.Println(err)
	}
	re.TaskID = taskID

	log.Println("Storage GetTaskStatus - bye")

	return re, nil
}
