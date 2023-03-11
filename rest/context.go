package rest

import (
	"net/http"

	"github.com/gopalrohra/gofly/transformers"
)

type FlyAPIContext struct {
	user     User
	request  *http.Request
	response FlyAPIResponse
	rt       transformers.RequestTransformer
}

func NewFlyAPIContext(r *http.Request, rt transformers.RequestTransformer) *FlyAPIContext {
	return &FlyAPIContext{request: r, rt: rt}
}
func (c *FlyAPIContext) Bind(i interface{}) {
	c.rt.PopulateData(i)
}
func (c *FlyAPIContext) IsAuthenticated() bool {
	return c.user.UserID > 0
}
func (c *FlyAPIContext) GetUser() User {
	return c.user
}
func (c *FlyAPIContext) SetUser(user User) {
	c.user = user
}
func (c *FlyAPIContext) AuthHeader() string {
	h := c.request.Header["Authorization"]
	if len(h) > 0 {
		return h[0]
	} else {
		return ""
	}
}
func (c *FlyAPIContext) JSON(res FlyAPIResponse) {
	c.response = res
}
func (c *FlyAPIContext) GetResponse() FlyAPIResponse {
	return c.response
}
