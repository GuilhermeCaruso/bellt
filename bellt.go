// Copyright 2019 Guilherme Caruso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bellt

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	//  Methods used for validation and comparison of the HandleFunc and
	//  Modules functions.
	methods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
	}
	mainRouter *Router

	// HTTPStatusCode verbose status code used for quick debugging Http.
	HTTPStatusCode = map[string]int{

		"CONTINUE":                        100,
		"SWITCHING_PROTOCOLS":             101,
		"OK":                              200,
		"CREATED":                         201,
		"ACCEPTED":                        202,
		"NON_AUTHORITATIVE_INFORMATION":   203,
		"NO_CONTENT":                      204,
		"RESET_CONTENT":                   205,
		"PARTIAL_CONTENT":                 206,
		"MULTIPLE_CHOICES":                300,
		"AMBIGUOUS":                       300,
		"MOVED_PERMANENTLY":               301,
		"MOVED":                           301,
		"FOUND":                           302,
		"REDIRECT":                        302,
		"SEE_OTHER":                       303,
		"REDIRECT_METHOD":                 303,
		"NOT_MODIFIED":                    304,
		"USE_PROXY":                       305,
		"UNUSED":                          306,
		"TEMPORARY_REDIRECT":              307,
		"REDIRECT_KEEP_VERB":              307,
		"BAD_REQUEST":                     400,
		"UNAUTHORIZED":                    401,
		"PAYMENT_REQUIRED":                402,
		"FORBIDDEN":                       403,
		"NOT_FOUND":                       404,
		"METHOD_NOT_ALLOWED":              405,
		"NOT_ACCEPTABLE":                  406,
		"PROXY_AUTHENTICATION_REQUIRED":   407,
		"REQUEST_TIMEOUT":                 408,
		"CONFLICT":                        409,
		"GONE":                            410,
		"LENGTH_REQUIRED":                 411,
		"PRECONDITION_FAILED":             412,
		"REQUEST_ENTITY_TOOLARGE":         413,
		"REQUEST_URI_TOO_LONG":            414,
		"UNSUPPORTED_MEDIA_TYPE":          415,
		"REQUESTED_RANGE_NOT_SATISFIABLE": 416,
		"EXPECTATION_FAILED":              417,
		"UNPROCESSABLE_ENTITY":            422,
		"UPGRADE_REQUIRED":                426,
		"INTERNAL_SERVER_ERROR":           500,
		"NOT_IMPLEMENTED":                 501,
		"BAD_GATEWAY":                     502,
		"SERVICE_UNAVAILABLE":             503,
		"GATEWAY_TIMEOUT":                 504,
		"HTTP_VERSION_NOT_SUPPORTED":      505,
		"NETWORK_CONNECT_TIMEOUT":         599,
	}
)

// Router is a struct responsible for storing routes already available (Route)
// or routes that will still be available (BuiltRoute).
type Router struct {
	routes []*Route
	built  []*BuiltRoute
}

// Route is a struct responsible for storing basic information of a Route, with
// all its variable parameters recorded.
type Route struct {
	Path    string
	Handler http.HandlerFunc
	Params  []Variable
}

// SubHandle is a struct similar to Route, however its behavior must be related
// to GroupHandle, having all its behavior mirrored from a Route.
type SubHandle struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

// BuiltRoute is an internal pattern struct for routes that will be built at
// run time.
type BuiltRoute struct {
	TempPath string
	Handler  http.HandlerFunc
	Var      map[int]Variable
	KeyRoute string
	Methods  []string
}

// Variable is a struct that guarantees the correct mapping of variables used
// in built routes.
type Variable struct {
	Name  string
	Value string
}

// ParamReceiver is responsible to return params set on context
type ParamReceiver struct {
	request *http.Request
}

// Middleware is a type responsible for characterizing middleware functions
// that should be used in conjunction with bellt.Use().
type Middleware func(http.HandlerFunc) http.HandlerFunc

// NewRouter is responsible to initialize a "singleton" router instance.
func NewRouter() *Router {
	if mainRouter == nil {
		http.HandleFunc("/health", healthApplication)
		http.Handle("/", verifyBuiltRoutes(redirectBuiltRoute))
		mainRouter = &Router{}
	}
	return mainRouter
}

/*
	Router is a struct responsible for storing routes already available (Route)
	or routes that will still be available (BuiltRoute).

	Its initialization is done through the method NewRouter:

		router: = bellt.NewRouter ()

		func main () {
			[...]
			log.Fatal (http.ListenAndServe (": 8080", nil))
		}
*/

// Method to obtain router for interanl processing
func getRouter() *Router {
	return mainRouter
}

// RedirectBuiltRoute Performs code analysis assigning values to variables
// in execution time.
func redirectBuiltRoute(w http.ResponseWriter, r *http.Request) {
	selectedBuilt, params := getRequestParams(r.URL.Path)

	router := getRouter()
	for idx, varParam := range selectedBuilt.Var {
		selectedBuilt.Var[idx] = Variable{
			Name:  varParam.Name,
			Value: params[idx],
		}
	}
	var allParams []Variable
	for _, param := range selectedBuilt.Var {
		allParams = append(allParams, param)
	}
	router.createBuiltRoute(
		selectedBuilt.TempPath,
		selectedBuilt.Handler,
		selectedBuilt.Methods,
		selectedBuilt.Var)

	setRouteParams(gateMethod(
		selectedBuilt.Handler,
		selectedBuilt.Methods...),
		allParams).ServeHTTP(w, r)
}

// Use becomes responsible for executing all middlewares passed through a
// cascade method.
func Use(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {

	for x := len(middleware) - 1; x >= 0; x-- {
		mid := middleware[x]
		handler = mid(handler)
	}

	return handler
}

/*
	The Use application should be done as follows within the
	router.HandleFunc() method:

		router: = bellt.NewRouter ()

		func main () {
			router.HandleFunc("path/", bellt.Use(
				handlerFunc,
				middleware1,
				middleware2,
				...,
			), "GET")
			log.Fatal (http.ListenAndServe (": 8080", nil))
		}
*/

// ----------------------------------------------------------------------------
// Router methods
// ----------------------------------------------------------------------------

// HandleFunc function responsible for initializing a common route or built
// through the Router. All non-grouped routes must be initialized by this
// method.
func (r *Router) HandleFunc(path string, handleFunc http.HandlerFunc,
	methods ...string) {
	key, values := getBuiltRouteParams(path)
	if values != nil {
		valuesList := make(map[int]Variable)

		for idx, name := range values {
			valuesList[idx] = Variable{
				Name:  name[1],
				Value: "",
			}
		}

		builtRoute := &BuiltRoute{
			TempPath: path,
			Handler:  handleFunc,
			Var:      valuesList,
			KeyRoute: key,
			Methods:  methods,
		}

		r.built = append(r.built, builtRoute)

	} else {

		route := &Route{
			Path:    path,
			Handler: handleFunc,
		}
		r.routes = append(r.routes, route)

		route.methods(methods...)
	}

}

/*
	The HandleFunc application should be done as follows within the
	router.HandleFunc() method:

		router: = bellt.NewRouter ()

		func main () {
			router.HandleFunc("path/", bellt.Use(
				handlerFunc,
				middleware1,
				middleware2,
				...,
			), "GET")
			log.Fatal (http.ListenAndServe (": 8080", nil))
		}
*/

// HandleGroup used to create and define a group of sub-routes
func (r *Router) HandleGroup(mainPath string, sr ...*SubHandle) {
	for _, route := range sr {
		var buf bytes.Buffer
		buf.WriteString(mainPath)
		buf.WriteString(route.Path)
		r.HandleFunc(buf.String(), route.Handler, route.Methods...)
	}
}

/*
	The HandleGroup and SubHandleFunc application should be done as follows
	within the router.HandleFunc() method:

		router: = bellt.NewRouter ()

		func main () {

			router.HandleGroup("/api",

				router.SubHandleFunc("/oi/{id}", bellt.Use(
					healthApplication,
					middleware,
				), "GET", "PUT"),

				router.SubHandleFunc("/oi/{user}", bellt.Use(
					healthApplication,
					middleware,
				), "GET", "PUT"),

			)
			log.Fatal (http.ListenAndServe (": 8080", nil))
		}
*/

// SubHandleFunc is responsible for initializing a common or built route. Its
// use must be made within the scope of the HandleGroup() method, where the
// main path will be declared.
func (r *Router) SubHandleFunc(path string, handleFunc http.HandlerFunc,
	methods ...string) *SubHandle {

	handleDetail := &SubHandle{
		Handler: handleFunc,
		Path:    path,
		Methods: methods,
	}
	return handleDetail
}

// Internal method of route construction based on parameters passed in the
// HandleFunc, guaranteeing a valid and functional route.
func (r *Router) routeBuilder(path string, handleFunc http.HandlerFunc,
	params ...Variable) *Route {
	route := &Route{
		Handler: handleFunc,
		Path:    path,
		Params:  params,
	}

	r.routes = append(r.routes, route)
	return route
}

// Internal method responsible for standardizing built routes in order to
// generate valid models of used.
func (r *Router) createBuiltRoute(path string, handler http.HandlerFunc,
	methods []string, params map[int]Variable) {
	var (
		builtPath = path
		allParams []Variable
	)

	for _, param := range params {
		builtPath = strings.Replace(builtPath, "{"+param.Name+"}",
			param.Value, -1)
		allParams = append(allParams, param)
	}

	r.routeBuilder(builtPath, handler, allParams...).methods(methods...)
}

// ----------------------------------------------------------------------------
// Route methods
// ----------------------------------------------------------------------------

// Internal method responsible for validating if the request method used exists
// for the route presented.
func (r *Route) methods(methods ...string) {
	for _, method := range methods {
		if !checkMethod(method) {
			log.Fatal(fmt.Sprintf("Method %s on %s not allowed",
				method, r.Path))
		}
	}

	if len(r.Params) > 0 {
		http.HandleFunc(r.Path, headerBuilder(
			setRouteParams(gateMethod(r.Handler, methods...), r.Params)))
	} else {
		http.HandleFunc(r.Path, headerBuilder(gateMethod(r.Handler,
			methods...)))
	}

}

// Internal method that validates whether the value passed in methods matches
// existing values.
func checkMethod(m string) bool {
	for _, method := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// Router middlewares
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

		w.WriteHeader(HTTPStatusCode["NOT_FOUND"])
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

//	Method to obtain route params in a built route
func getBuiltRouteParams(path string) (string, [][]string) {
	rgx := regexp.MustCompile(`(?m){(\w*)}`)
	rgxStart := regexp.MustCompile(`(?m)(^\/)`)
	rgxEnd := regexp.MustCompile(`(?m)(\/$)`)
	return rgxEnd.ReplaceAllString(rgxStart.ReplaceAllString(
		rgx.Split(path, -1)[0], ""), ""), rgx.FindAllStringSubmatch(path, -1)
}

// Method to obtain request methods
func getRequestParams(path string) (*BuiltRoute, map[int]string) {
	router := getRouter()

	var builtRouteList *BuiltRoute
	params := make(map[int]string)

	for _, route := range router.built {
		rgx := regexp.MustCompile(route.KeyRoute)
		if rgx.FindString(path) != "" {
			if (len(strings.Split(
				rgx.Split(path, -1)[1], "/")) - 1) == len(route.Var) {
				builtRouteList = route
				for idx, val := range strings.Split(rgx.Split(path, -1)[1],
					"/") {
					if idx != 0 {
						params[idx-1] = val
					}
				}
			}
		}
	}
	return builtRouteList, params
}

// RouteVariables used to capture and store parameters passed to built routes
func RouteVariables(r *http.Request) *ParamReceiver {

	receiver := ParamReceiver{
		request: r,
	}

	return &receiver
}

// Defines and organizes route parameters by applying them in request
func setRouteParams(next http.HandlerFunc, params []Variable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, param := range params {
			ctx = context.WithValue(ctx, param.Name, param.Value)
		}

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// ----------------------------------------------------------------------------
// ParamReceiver middlewares
// ----------------------------------------------------------------------------

// GetVar return a value of router variable
func (pr *ParamReceiver) GetVar(variable string) interface{} {
	return pr.request.Context().Value(variable)
}

// ----------------------------------------------------------------------------
// Server support methods
// ----------------------------------------------------------------------------

// Function used in application health routing.
func healthApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusCode["OK"])
	w.Write([]byte(`{"alive": "Server running"}`))
}

// Converts and prepares Handle function for built route analysis
func verifyBuiltRoutes(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
