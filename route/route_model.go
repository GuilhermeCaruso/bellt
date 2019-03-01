package route

import "net/http"

// Router is a struct to define a router instance.
type Router struct {
	routes []*Route
	built  []*BuiltRoute
}

// Route is a struct to define a route item.
type Route struct {
	Path    string
	Handler http.HandlerFunc
}

// SubRoute is a struct to define a subRoute item.
type SubRoute struct {
	Route   Route
	Methods []string
}

// BuiltRoute is a struct to define a BuiltRoute item.
type BuiltRoute struct {
	TempPath string
	Handler  http.HandlerFunc
	Var      map[string]string
}
