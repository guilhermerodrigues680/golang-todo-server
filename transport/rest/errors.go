package rest

import "errors"

var ErrDecodeRequestBody = errors.New("failed to decode the request body")
var ErrInvalidRequestBody = errors.New("invalid request body")
var ErrUnknownService = errors.New("service error")
var ErrAuthUsernamePassword = errors.New("The user name or password is incorrect")
