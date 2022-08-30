package rest

type FlyAPIResource struct {
	Authenticator FlyAPIHandler
	AllowedRoles  []string
	NewController ControllerProvider
	Handler       FlyAPIHandler
}
type FlyAPIHandler = func(ctx *FlyAPIContext)
type ControllerProvider = func(*FlyAPIContext) FlyAPIController

type User struct {
	UserID int
	Role   string
	Email  string
}
