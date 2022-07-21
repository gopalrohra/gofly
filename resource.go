package flyapi

type Bind = func(interface{})
type FlyAPIResource struct {
	Authenticate bool
	AllowedRoles []string
	Controller   FlyAPIController
}
type FlyAPIController interface {
	Init(FlyAPIContext, Bind) FlyAPIController
	Parse() FlyAPIController
	Authorize() FlyAPIController
	Execute() FlyAPIController
	GetResponse() FlyAPIResponse
}
type User struct {
	UserID int
	Role   string
	Email  string
}
type FlyAPIContext struct {
	User  User
	Paths []string
	Data  interface{}
}
type FlyAPIResponse struct {
	Status  string
	Message string
	Data    interface{}
}
