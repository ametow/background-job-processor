package storage

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/ametow/background-job-processor/internal/server/domain/entity"
)

func TestTaskStorage_CreateTask_RealDB(t *testing.T) {
	// TODO make tests independent of main code
	ts := NewTaskStorage()

	type args struct {
		te entity.TaskEntity
	}
	tests := []struct {
		name    string
		fields  *TaskStorage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "created",
			fields: ts,
			args: args{
				te: entity.TaskEntity{
					ID:      "abcde",
					Method:  http.MethodGet,
					URL:     "http://google.com",
					Headers: entity.Headers{"Authentication": "Basic bG9naW46cGFzc3dvcmQ=", "Content-type": "JSON"},
				},
			},
			wantErr: false,
		},
		{
			name:   "created, this task is to check its own status",
			fields: ts,
			args: args{
				te: entity.TaskEntity{
					ID:      "checkstatus",
					Method:  http.MethodGet,
					URL:     "http://localhost:8080/task/checkstatus",
					Headers: entity.Headers{"Content-type": "JSON"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TaskStorage{
				DB: tt.fields.DB,
			}
			if err := ts.CreateTask(tt.args.te); (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskStorage_Create30Tasks_RealDB(t *testing.T) {

	ts := NewTaskStorage()

	type args struct {
		te entity.TaskEntity
	}
	tests := []struct {
		name    string
		fields  *TaskStorage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "created",
			fields: ts,
			args: args{
				te: entity.TaskEntity{
					ID:      "boredapi",
					Method:  http.MethodGet,
					URL:     "https://www.boredapi.com/api/activity",
					Headers: entity.Headers{"Content-type": "JSON"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TaskStorage{
				DB: tt.fields.DB,
			}

			for i := 0; i < 30; i++ {
				tt.args.te.ID = "testtask" + fmt.Sprint(i)
				if err := ts.CreateTask(tt.args.te); (err != nil) != tt.wantErr {
					t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			// TODO run agent
			// TODO clean Data

		})
	}
}

// TODO add DB clean-up - garbage in, garbage out
func TestTaskStorage_GetTaskStatus_RealDB(t *testing.T) {

	ts := NewTaskStorage()

	type args struct {
		taskID string
		result entity.ResultEntity
	}
	tests := []struct {
		name    string
		fields  *TaskStorage
		args    args
		want    entity.ResultEntity
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "found",
			args: args{
				taskID: "abc",
				result: entity.ResultEntity{
					TaskID:                 "abc",
					TaskStatus:             "done",
					ResponseHttpStatusCode: 200,
					ResponseHeaders:        entity.Headers{"Result-Header1": "sample-data", "Content-type": "JSON"},
					ResponseBodyLength:     "50",
				},
			},
			fields: ts,
			want: entity.ResultEntity{
				TaskID:                 "abc",
				TaskStatus:             "done",
				ResponseHttpStatusCode: 200,
				ResponseHeaders:        entity.Headers{"Result-Header1": "sample-data", "Content-type": "JSON"},
				ResponseBodyLength:     "50",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TaskStorage{
				DB: tt.fields.DB,
			}

			_, err := ts.DB.Exec(
				`INSERT INTO tasks (task_id, task_status, result_http_status_code, result_headers, result_body_length)
				VALUES ($1, $2, $3, $4, $5);`,
				tt.args.result.TaskID, tt.args.result.TaskStatus, tt.args.result.ResponseHttpStatusCode, tt.args.result.ResponseHeaders, tt.args.result.ResponseBodyLength,
			)
			if err != nil {
				log.Fatal(err)
			}

			got, err := ts.GetTaskStatus(tt.args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
