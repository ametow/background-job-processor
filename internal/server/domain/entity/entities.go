package entity

type Headers map[string]string

type TaskEntity struct {
	ID      string
	Method  string
	URL     string
	Headers Headers
}

type ResultEntity struct {
	TaskID                 string  `json:"id"`
	TaskStatus             string  `json:"status"`
	ResponseHttpStatusCode int     `json:"httpStatusCode"`
	ResponseHeaders        Headers `json:"headers"`
	ResponseBodyLength     string  `json:"length"`
}
