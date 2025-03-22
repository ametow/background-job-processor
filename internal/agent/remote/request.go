package remote

import (
	"log"
	"net/http"

	"github.com/ametow/background-job-processor/internal/server/domain/entity"
)

type Request struct{}

func NewRequest() Request {
	return Request{}
}

func (r Request) Request(taskEntity entity.TaskEntity) (entity.ResultEntity, error) {
	log.Println("Remote.Request - start")

	result := entity.ResultEntity{}
	client := &http.Client{}

	request, err := http.NewRequest(taskEntity.Method, taskEntity.URL, nil)
	if err != nil {
		log.Println(err)
		return entity.ResultEntity{}, err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		result.TaskStatus = entity.ERROR
		return result, err
	}

	headers := make(entity.Headers)
	for k, v := range response.Header {
		headers[k] = v[0]
	}

	result.TaskID = taskEntity.ID
	result.TaskStatus = entity.DONE
	result.ResponseHttpStatusCode = response.StatusCode
	result.ResponseHeaders = headers
	result.ResponseBodyLength = response.Header.Get("Content-Length")

	log.Println("Remote.Request - Result is:", result)

	log.Println("Remote.Request - end")
	return result, nil
}
