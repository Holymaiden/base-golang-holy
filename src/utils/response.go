package utils

import (
	"jwt/src/requests"
)

type EmptyObj struct{}

func Response(code bool, message string, data interface{}, errors error) requests.Response {
	return requests.Response{
		Code:    code,
		Message: message,
		Data:    data,
		Errors:  errors,
	}
}

func ResponseError(message string, data interface{}, errors error) requests.Response {
	return Response(false, message, data, errors)
}
