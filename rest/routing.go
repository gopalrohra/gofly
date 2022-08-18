package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gopalrohra/flyapi/transformers"
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
	fmt.Printf("Value of route path: %s and value of request path: %s\n", route.Path, r.URL.Path)
	processRoute(w, r, route)
}
func (router *Router) match(path string, method string) Route {
	paths := strings.Split(path, "/")
	for _, route := range router.Routes {
		if path == route.Path {
			return route
		}
		routePaths := strings.Split(route.Path, "/")
		if len(paths) == len(routePaths) && elementsAreSame(paths, routePaths) {
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
		t := transformers.RequestTransformer{Request: r, RoutePath: route.Path}
		t.ParseParameters()
		response := processController(resource.NewController(), ctx, t)
		fmt.Fprint(w, util.ToJSONString(response))
	} else {
		fmt.Fprint(w, util.ToJSONString(MethodNotAllowedResponse))
	}
}
func processController(controller FlyAPIController, ctx FlyAPIContext, t transformers.RequestTransformer) FlyAPIResponse {
	controller.Init(ctx, t.PopulateData)
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