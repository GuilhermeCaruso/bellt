package bellt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestRedirectBuiltRoute(t *testing.T) {
	req, err := http.NewRequest("GET", "/usuario/123/novo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.Handler(verifyBuiltRoutes(redirectBuiltRoute))

	handler.ServeHTTP(rr, req)

	fmt.Println(rr, req)
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// expected := `{"alive": "Server running"}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
