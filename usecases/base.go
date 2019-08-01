package usecases

import (

)

// UseCase is an interface which standardises the methods available on
// Use Cases.
type UseCase interface {
	Setup()
	Execute()
}

type ResponseError struct {
	Name string
	Description string
}

type Response struct {
	Body string
}

type ResponseCollector struct {
	Response *Response
	Error *ResponseError
}

func (rc *ResponseCollector) SetResponse(resp *Response) {
	rc.Response = resp
	rc.Error = nil
}

func (rc *ResponseCollector) SetError(err *ResponseError) {
	rc.Response = nil
	rc.Error = err
}

func panicHandler(response *ResponseCollector) {
	if err := recover(); err != nil {
		respErr := ResponseError{
			Name: "SEVERE_FAILURE",
			Description: "Something went terribly wrong.",
		}
		response.SetError(&respErr)
	}
}