package rest

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
