package route

import (
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

// ----------------------------------------------------------------------------
// Route methods
// ----------------------------------------------------------------------------

//Method define request handler method
func (r *Route) Method(method string) {
	if checkMethod(method) {
		http.HandleFunc(r.Path, headerBuilder(gateMethod(method, r.Handler)))
		return
	}
	log.Fatal()
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
func gateMethod(method string, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(utils.HTTPStatusCode["NOT_FOUND"])
		w.Write([]byte(`{"error": "The method for this route does not exist"}`))

	}
}

//  Defines JSON header for standard REST service routes.
func headerBuilder(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}
