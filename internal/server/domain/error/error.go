package error

import "errors"

var NotFound = errors.New("not found")
var RecordExists = errors.New("record exists")
var WrongContentType = errors.New("wrong content type in request")
