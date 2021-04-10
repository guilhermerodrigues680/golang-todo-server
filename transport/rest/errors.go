package rest

import "errors"

var ErrDecodeRequestBody = errors.New("Error decoding request body")
var ErrInvalidRequestBody = errors.New("Error invalid request body")
var ErrUnknownService = errors.New("Service error")
