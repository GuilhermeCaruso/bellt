package bellt

import (
	"net/http"
	"testing"
)

type simpleRouteTest struct {
	path    string
	handler http.HandlerFunc
}

func testHandler(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	globalRouter := getRouter()
	if router != globalRouter {
		t.Errorf("Wrong router initialization.\n\tgot:%+v\n\twant:%+v", router, globalRouter)
	}
}

func TesteRouteCreation(t *testing.T) {
	router := NewRouter()
	simpleRoutes := []simpleRouteTest{
		{
			path:    "/routeOne",
			handler: testHandler,
		},
		{
			path:    "/routeTwo",
			handler: testHandler,
		},
		{
			path:    "/routeThree",
			handler: testHandler,
		},
	}

	for _, route := range simpleRoutes {
		router.Route(route.path, route.handler)
	}

	if len(router.routes) != len(simpleRoutes) {
		t.Errorf("Invalid route number .\n\tgot:%+v\n\twant:%+v", len(router.routes), len(simpleRoutes))
	}

}
