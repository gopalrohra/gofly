package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gopalrohra/gofly/log"
	"github.com/gopalrohra/gofly/transformers"
	"github.com/gopalrohra/gofly/util"
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

func (router *Router) HandleRouting(w http.ResponseWriter, r *http.Request) {
	log.Infof("Processing %v for %v\n", r.Method, r.URL)
	route := router.match(r.URL.Path, r.Method)
	if route.Path == "" {
		fmt.Fprintf(w, util.ToJSONString(NotFoundResponse))
		return
	}
	log.Debugf("Value of route path: %s and value of request path: %s\n", route.Path, r.URL.Path)
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
		ctx := NewFlyAPIContext(r, transformers.NewRequestTransformer(r, route.Path))
		if resource.Authenticator != nil {
			resource.Authenticator(ctx)
			if !ctx.IsAuthenticated() {
				fmt.Fprintf(w, util.ToJSONString(AuthenticationErrorResponse))
				return
			}
		}
		response := processResource(ctx, resource)
		fmt.Fprint(w, util.ToJSONString(response))
	} else {
		fmt.Fprint(w, util.ToJSONString(MethodNotAllowedResponse))
	}
}

func processResource(ctx *FlyAPIContext, resource FlyAPIResource) FlyAPIResponse {
	if resource.NewController != nil {
		return processController(resource.NewController(ctx))
	} else if resource.Handler != nil {
		resource.Handler(ctx)
		return ctx.GetResponse()
	} else {
		return NoHandlerFoundResponse
	}
}
func processController(controller FlyAPIController) FlyAPIResponse {
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
