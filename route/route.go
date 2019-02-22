package route

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt/utils"
)

// NewRouter is responsible to initialize a "singleton" router instance.
func NewRouter() *Router {
	router := Router{}
	return &router
}

// ----------------------------------------------------------------------------
// Router methods
// ----------------------------------------------------------------------------

// Route used to define the calling method.
func (r *Router) Route(path string, handleFunc http.HandlerFunc) *Route {
	route := &Route{
		Path:    path,
		Handler: handleFunc,
	}
	r.routes = append(r.routes, route)

	return route
}

// SubRoute used to create and define a sub-route belonging to a route grouping.
func (r *Router) SubRoute(path string, handleFunc http.HandlerFunc, methods ...string) *SubRoute {

	subRoute := &SubRoute{
		Route: Route{
			Path:    path,
			Handler: handleFunc,
		},
		Methods: methods,
	}

	return subRoute
}

// Group used to create and define a group of sub-routes
func (r *Router) Group(mainPath string, sr ...*SubRoute) {

	for _, route := range sr {
		var buf bytes.Buffer
		buf.WriteString(mainPath)
		buf.WriteString(route.Route.Path)
		r.Route(buf.String(), route.Route.Handler).Methods(route.Methods...)
	}
}

// ----------------------------------------------------------------------------
// Route methods
// ----------------------------------------------------------------------------

//Methods define request handler method
func (r *Route) Methods(methods ...string) {
	for _, method := range methods {
		if !checkMethod(method) {
			log.Fatal(fmt.Sprintf("Method %s on %s not allowed", method, r.Path))
		}
	}
	http.HandleFunc(r.Path, headerBuilder(gateMethod(r.Handler, methods...)))

}

// Verifies veracity of an established method.
func checkMethod(m string) bool {
	for _, method := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// Route middlewares
// ----------------------------------------------------------------------------

// Ensures that routing is done using valid methods
func gateMethod(next http.HandlerFunc, methods ...string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(utils.HTTPStatusCode["NOT_FOUND"])
		w.Write([]byte(`{"error": "The method for this route doesnt exist"}`))

	}
}

// Defines JSON header for standard REST service routes.
func headerBuilder(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}
