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

// Route method called to create a new simple route on Router.
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

func (r *Route) Method(method string) {
	if checkMethod(method) {
		http.HandleFunc(r.Path, headerBuilder(gateMethod(method, r.Handler)))
		return
	}
	log.Fatal()
}

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

func gateMethod(method string, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			next.ServeHTTP(w, r)
		}

		w.WriteHeader(utils.HTTPStatusCode["NOT_FOUND"])
		w.Write([]byte(`{"error": "The method for this route does not exist"}`))

	}
}

func headerBuilder(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}
