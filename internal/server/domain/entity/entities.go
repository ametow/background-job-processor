package entity

// http headers
type Headers map[string]string

type TaskStatus string

const (
	DONE       TaskStatus = "done"
	ERROR      TaskStatus = "error"
	NEW        TaskStatus = "new"
	IN_PROCESS TaskStatus = "in_process"
)

type TaskEntity struct {
	ID      string
	Method  string
	URL     string
	Headers Headers
}

type ResultEntity struct {
	TaskID                 string     `json:"id"`
	TaskStatus             TaskStatus `json:"status"`
	ResponseHttpStatusCode int        `json:"httpStatusCode"`
	ResponseHeaders        Headers    `json:"headers"`
	ResponseBodyLength     string     `json:"length"`
}
