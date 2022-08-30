package rest

import "fmt"

type FlyAPIController interface {
	Validate()
	Authorize()
	Execute()
	HasErrors() bool
	GetResponse() FlyAPIResponse
}
type UnImplementedFlyController struct{}

func (c *UnImplementedFlyController) Validate() {
	fmt.Println("Validate method not implemented")
}
func (c *UnImplementedFlyController) Authorize() {
	fmt.Println("Authorize method not implemented")
}
func (c *UnImplementedFlyController) Execute() {
	fmt.Println("Execute method not implemented")
}
func (c *UnImplementedFlyController) GetResponse() FlyAPIResponse {
	fmt.Println("GetResponse method not implemented")
	return FlyAPIResponse{Status: InternalServerError, Message: InternalServerErrorMessage}
}
