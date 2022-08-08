package rest

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	var routes []Route
	for i := 0; i < 5; i++ {
		route := Route{
			Path: fmt.Sprintf("/orders%v/:id", i+1),
			Resources: map[string]FlyAPIResource{
				"GET": {},
			},
		}
		routes = append(routes, route)
	}
	router := Router{Routes: routes}
	path := "/orders1/3"
	routeResult := router.match(path, "GET")
	if "/orders1/:id" != routeResult.Path {
		t.Errorf("Expected route: %s and got route: %s\n", path, routeResult.Path)
	}
}
