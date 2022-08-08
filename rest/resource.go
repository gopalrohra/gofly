package rest

import "net/http"

type Bind = func(interface{})
type AuthFunc = func(w http.ResponseWriter, r *http.Request) FlyAPIContext
type FlyAPIResource struct {
	Authenticator AuthFunc
	AllowedRoles  []string
	NewController ControllerProvider
}
type ControllerProvider = func() FlyAPIController

type User struct {
	UserID int
	Role   string
	Email  string
}
type FlyAPIContext struct {
	User User
	Data interface{}
}
type FlyAPIResponse struct {
	Status  string
	Message string
	Data    interface{}
}
