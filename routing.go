package flyapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gopalrohra/flyapi/util"
	"github.com/rs/cors"
)

type Route struct {
	Path           string
	AllowedMethods STRList
	Authenticate   bool
	AllowedRoles   []string
	Resource       FlyAPIResource
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
	//flyAPIContext := r.Context().Value("FlyAPIContext").(FlyAPIContext)
	flyAPIContext := FlyAPIContext{User: User{UserID: 1, Email: "gopal.rohra@gmail.com"}}
	route := router.match(r.URL.Path, r.Method)
	//todo: if route is not found return 404
	fmt.Printf("Allowed methods on path: %s are %s\n", route.Path, route.AllowedMethods)
	fmt.Println("Inside allowedMethods condition")
	//fmt.Printf("Kind of data: %v\n", reflect.ValueOf(data).Elem())
	t := RequestTransformer{request: r}
	resource := route.Resource
	response := resource.Init(flyAPIContext, t.populateData).Parse().Authorize().Execute().GetResponse()
	fmt.Fprint(w, util.ToJSONString(response))
	return
	//fmt.Fprint(w, util.ToJSONString(methodNotAllowedResponse))
}
func (router *Router) match(path string, method string) Route {
	paths := strings.Split(path, "/")
	for _, route := range router.Routes {
		if path == route.Path && route.AllowedMethods.Contains(method) {
			return route
		}
		routePaths := strings.Split(route.Path, "/")
		if len(paths) == len(routePaths) && elementsAreSame(paths, routePaths) && route.AllowedMethods.Contains(method) {
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
