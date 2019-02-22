package route

import "net/http"

// Router is a struct to define a router instance.
type Router struct {
	routes []*Route
}

// Route is a struct to define a route item.
type Route struct {
	Path    string
	Handler http.HandlerFunc
}
