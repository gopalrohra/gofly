package flyapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gopalrohra/flyapi/util"
	"github.com/rs/cors"
)

type Route struct {
	Path      string
	Resources map[string]FlyAPIResource
}

type Router struct {
	Routes []Route
	Cors   *cors.Cors
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

func (router *Router) HandleRouting(w http.ResponseWriter, r *http.Request) {
	route := router.match(r.URL.Path, r.Method)
	if route.Path == "" {
		fmt.Fprintf(w, util.ToJSONString(NotFoundResponse))
		return
	}
	processRoute(w, r, route)
}
func (router *Router) match(path string, method string) Route {
	paths := strings.Split(path, "/")
	for _, route := range router.Routes {
		if path == route.Path { //&& route.AllowedMethods.Contains(method) {
			return route
		}
		routePaths := strings.Split(route.Path, "/")
		if len(paths) == len(routePaths) && elementsAreSame(paths, routePaths) { //&& route.AllowedMethods.Contains(method) {
			return route
		}
	}
	return Route{}
}
func elementsAreSame(paths []string, routePaths []string) bool {
	for i, path := range paths {
		if path != routePaths[i] && !strings.HasPrefix(routePaths[i], ":") {
			return false
		}
	}
	return true
}
func processRoute(w http.ResponseWriter, r *http.Request, route Route) {
	if resource, ok := route.Resources[r.Method]; ok {
		var ctx FlyAPIContext
		if resource.Authenticator != nil {
			ctx = resource.Authenticator(w, r)
			if ctx.User.UserID == 0 {
				fmt.Fprintf(w, util.ToJSONString(AuthenticationErrorResponse))
				return
			}
		}
		t := RequestTransformer{request: r}
		t.parseParameters()
		response := processController(resource.NewController(), ctx, t)
		fmt.Fprint(w, util.ToJSONString(response))
	} else {
		fmt.Fprint(w, util.ToJSONString(MethodNotAllowedResponse))
	}
}
func processController(controller FlyAPIController, ctx FlyAPIContext, t RequestTransformer) FlyAPIResponse {
	controller.Init(ctx, t.populateData)
	if !controller.HasErrors() {
		controller.Validate()
	}
	if !controller.HasErrors() {
		controller.Authorize()
	}
	if !controller.HasErrors() {
		controller.Execute()
	}
	return controller.GetResponse()
}
