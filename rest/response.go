package rest

import "github.com/gopalrohra/gofly/log"

type FlyAPIResponse struct {
	Status  string
	Message string
	Data    interface{}
}

var MethodNotAllowedResponse = FlyAPIResponse{
	Status:  "MethodNotAllowed",
	Message: "This method is not allowed on this resource",
}
var NotFoundResponse = FlyAPIResponse{
	Status:  "NotFound",
	Message: "Resource not found",
}
var AuthenticationErrorResponse = FlyAPIResponse{
	Status:  "AuthenticationError",
	Message: "Can't verify the identity",
}
var NoHandlerFoundResponse = FlyAPIResponse{
	Status:  InternalServerError,
	Message: InternalServerErrorMessage,
}

func BuildBadRequestResponse(validationErrors []string) FlyAPIResponse {
	return FlyAPIResponse{
		Status:  BadRequest,
		Message: BadRequestMessage,
		Data:    validationErrors,
	}
}
func BuildInternalServerResponse(logMsg string, err error) FlyAPIResponse {
	log.Errorf("%v: %v", logMsg, err)
	return FlyAPIResponse{
		Status:  InternalServerError,
		Message: InternalServerErrorMessage,
	}
}
func BuildSuccessResponse(entity interface{}) FlyAPIResponse {
	return FlyAPIResponse{
		Status:  "Success",
		Message: "Record",
		Data:    entity,
	}
}
