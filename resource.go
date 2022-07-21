package flyapi

type Bind = func(interface{})
type FlyAPIResource interface {
	Init(FlyAPIContext, Bind) FlyAPIResource
	Parse() FlyAPIResource
	Authorize() FlyAPIResource
	Execute() FlyAPIResource
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
