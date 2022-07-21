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

var methodNotAllowedResponse = FlyAPIResponse{
	Status:  "MethodNotAllowed",
	Message: "This method is not allowed on this resource",
}

func (router *Router) HandleRouting(w http.ResponseWriter, r *http.Request) {
	route := router.match(r.URL.Path, r.Method)
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
	//flyAPIContext := r.Context().Value("FlyAPIContext").(FlyAPIContext)
	flyAPIContext := FlyAPIContext{User: User{UserID: 1, Email: "gopal.rohra@gmail.com"}}
	//todo: if route is not found return 404
	//fmt.Printf("Allowed methods on path: %s are %s\n", route.Path, route.AllowedMethods)
	//t := RequestTransformer{request: r}
	if resource, ok := route.Resources[r.Method]; ok {
		t := RequestTransformer{request: r}
		controller := resource.Controller
		response := controller.Init(flyAPIContext, t.populateData).Parse().Authorize().Execute().GetResponse()
		fmt.Fprint(w, util.ToJSONString(response))
	} else {
		fmt.Fprint(w, util.ToJSONString(methodNotAllowedResponse))
	}
	//resource := route.Resource
	//response := resource.Init(flyAPIContext, t.populateData).Parse().Authorize().Execute().GetResponse()
	//fmt.Fprint(w, util.ToJSONString(response))
}
