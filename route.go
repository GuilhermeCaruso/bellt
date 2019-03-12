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

//	Method call definition constants
var (
	methods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
	}
	mainRouter *Router
)

// Router is a struct to define a router instance.
type Router struct {
	routes []*Route
	built  []*BuiltRoute
}

// Route is a struct to define a route item.
type Route struct {
	Path    string
	Handler http.HandlerFunc
	Params  []Variable
}

// SubRoute is a struct to define a subRoute item.
type SubHandle struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

// BuiltRoute is a struct to define a BuiltRoute item.
type BuiltRoute struct {
	TempPath string
	Handler  http.HandlerFunc
	Var      map[int]Variable
	KeyRoute string
	Methods  []string
}

// Variable is a struct responsible to design a dinamic variable
type Variable struct {
	Name  string
	Value string
}

// ParamReceiver Responsible to return params set on context
type ParamReceiver struct {
	request *http.Request
}

// NewRouter is responsible to initialize a "singleton" router instance.
func NewRouter() *Router {
	if mainRouter == nil {
		mainRouter = &Router{}
	}
	return mainRouter
}

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
	router.createBuiltRoute(selectedBuilt.TempPath, selectedBuilt.Handler, selectedBuilt.Methods, selectedBuilt.Var)

	setRouteParams(gateMethod(selectedBuilt.Handler, selectedBuilt.Methods...), allParams).ServeHTTP(w, r)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Use(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {

	for x := len(middlewares) - 1; x >= 0; x-- {
		mid := middlewares[x]
		handler = mid(handler)
	}

	return handler
}

// ----------------------------------------------------------------------------
// Router methods
// ----------------------------------------------------------------------------

func (r *Router) SubHandleFunc(path string, handleFunc http.HandlerFunc, methods ...string) *SubHandle {

	handleDetail := &SubHandle{
		Handler: handleFunc,
		Path:    path,
		Methods: methods,
	}
	return handleDetail
}

func (r *Router) HandleFunc(path string, handleFunc http.HandlerFunc, methods ...string) {
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

// Used to define built call method
func (r *Router) routeBuilder(path string, handleFunc http.HandlerFunc, params ...Variable) *Route {
	route := &Route{
		Handler: handleFunc,
		Path:    path,
		Params:  params,
	}

	r.routes = append(r.routes, route)
	return route
}

// Creation and modeling of built route
func (r *Router) createBuiltRoute(path string, handler http.HandlerFunc, methods []string, params map[int]Variable) {
	var (
		builtPath = path
		allParams []Variable
	)

	for _, param := range params {
		builtPath = strings.Replace(builtPath, "{"+param.Name+"}", param.Value, -1)
		allParams = append(allParams, param)
	}

	r.routeBuilder(builtPath, handler, allParams...).methods(methods...)
}

// HandleGroup used to create and define a group of sub-routes
func (r *Router) HandleGroup(mainPath string, sr ...*SubHandle) {

	for _, route := range sr {
		var buf bytes.Buffer
		buf.WriteString(mainPath)
		buf.WriteString(route.Path)
		r.HandleFunc(buf.String(), route.Handler, route.Methods...)
	}
}

// ----------------------------------------------------------------------------
// Route methods
// ----------------------------------------------------------------------------

//Methods define request handler method
func (r *Route) methods(methods ...string) {
	for _, method := range methods {
		if !checkMethod(method) {
			log.Fatal(fmt.Sprintf("Method %s on %s not allowed", method, r.Path))
		}
	}

	if len(r.Params) > 0 {
		http.HandleFunc(r.Path, headerBuilder(setRouteParams(gateMethod(r.Handler, methods...), r.Params)))
	} else {
		http.HandleFunc(r.Path, headerBuilder(gateMethod(r.Handler, methods...)))
	}

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
	return rgxEnd.ReplaceAllString(rgxStart.ReplaceAllString(rgx.Split(path, -1)[0], ""), ""), rgx.FindAllStringSubmatch(path, -1)
}

// Method to obtain request methods
func getRequestParams(path string) (*BuiltRoute, map[int]string) {
	router := getRouter()

	var builtRouteList *BuiltRoute
	params := make(map[int]string)

	for _, route := range router.built {
		rgx := regexp.MustCompile(route.KeyRoute)
		if rgx.FindString(path) != "" {
			if (len(strings.Split(rgx.Split(path, -1)[1], "/")) - 1) == len(route.Var) {
				builtRouteList = route
				for idx, val := range strings.Split(rgx.Split(path, -1)[1], "/") {
					if idx != 0 {
						params[idx-1] = val
					}
				}
			}
		}
	}
	return builtRouteList, params
}

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
