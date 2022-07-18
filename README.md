# Bellt  
> Simple Golang HTTP router

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/GuilhermeCaruso/bellt)](https://goreportcard.com/report/github.com/GuilhermeCaruso/bellt) [![codecov](https://codecov.io/gh/GuilhermeCaruso/bellt/branch/master/graph/badge.svg)](https://codecov.io/gh/GuilhermeCaruso/bellt) [![Build Status](https://travis-ci.com/GuilhermeCaruso/bellt.svg?branch=master)](https://travis-ci.com/GuilhermeCaruso/bellt) ![GitHub](https://img.shields.io/badge/golang%20->=1.7-blue.svg) [![GoDoc](https://godoc.org/github.com/GuilhermeCaruso/bellt?status.svg)](https://godoc.org/github.com/GuilhermeCaruso/bellt) 

<p align="left">
    <img width="150" src="./assets/logo.png">
</p>

Bellt Package implements a request router with the aim of managing controller actions based on fixed and parameterized routes.

The project so far has the following functionalities:

* Standard definition of route "/health", in order to prepare the service developed with bellt to act as microservice.
* Providing the creation of parameterized routes, simple or segmented (groups).
* All requests can be made through fixed patterns, querystrings and parameters.
* Obtaining the requisition parameters in the controller functions.
---
# Summary
 * [Install](#install)
 * [Guide](#guide)
	* [Router](#router)
		* [HandleFunc](#handleFunc)
		* [HandleGroup](#handleGroup)
		* [SubHandleFunc](#subHandleFunc)
	* [Middleware](#middleware)
		* [Use](#use)
	* [Parameterized Routes](#parameterized-routes)
		* [Route Variables](#route-variables)
			* [GetVar](#getVar)
 * [Full Example](#full-example)
 * [Benchmark](#benchmark)
 * [Author](#author)
 * [Presentation](#presentation)
 * [License](#license)

# Install


To get Bellt

##### > Go CLI
```sh
go get -u github.com/GuilhermeCaruso/bellt
```
##### > Go DEP
```sh
dep ensure -add github.com/GuilhermeCaruso/bellt
```
##### > Govendor
```sh
govendor fetch github.com/GuilhermeCaruso/bellt
```

# Guide

## Router

To initialize our router
```go
var router = bellt.NewRouter()
```

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {
	router := bellt.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

### HandleFunc   

HandleFunc function responsible for initializing a common route or built through the Router. All non-grouped routes must be initialized by this method.

```go
/*
	[path] - Endpoint string
	[handlerFunc] - Function that will be called on the request
	[methods] - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
*/

router.HandleFunc(path, handlerFunc, methods)
    
```
```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {
	router := bellt.NewRouter()

	router.HandleFunc("/bellt", belltHandler, "GET")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func belltHandle(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simple Golang HTTP router")
}

```

### HandleGroup   

HandleGroup is responsible for creating a group of routes. The main path can be set for all other routes.

```go
/*
	[mainPath] - Main route used in all subr-outes
	
	[subHandleFunc] - SubHandleFunc function responsiblefor initializing a common route or
	built through the Router. All grouped routes must be initialized by this method
*/

router.HandleGroup(mainPath, ...SubHandleFunc)
    
```

### SubHandleFunc  

SubHandleFunc is responsible for initializing a common or built route. Its use must be made within the scope of the HandleGroup method, where the main path will be declared.

```go
/*
	[path] - Endpoint string
	[handlerFunc] - Function that will be called on the request
	[methods] - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
*/

router.SubHandleFunc(path, handlerFunc, methods)
    
```
```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {
	router := bellt.NewRouter()

	router.HandleGroup("/api",
		router.SubHandleFunc("/bellt", belltHandle, "GET"),
		router.SubHandleFunc("/check", checkHandle, "GET"),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func belltHandle(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simple Golang HTTP router")
}

func checkHandle(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok!")
}

```

## Middleware

The declaration of middlewares in HandleFunc or SubHandleFunc should be done using the *Use* method

### Use

```go
/*
	handlerFunc - Function that will be called on the request 
	middlewareList - Slice of middleware that will be used in the request (Middleware)
*/
bellt.Use(handlerFunc, ...middlewareList)
```

The middleware type has a following signature


```go
type Middleware func(http.HandlerFunc) http.HandlerFunc
```

Applying middlewares to routes

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {

	router := bellt.NewRouter()

	router.HandleFunc("/hello", bellt.Use(
		exampleHandler,
		middlewareOne,
		middlewareTwo,
	), "GET")

	router.HandleGroup("/api",
		router.SubHandleFunc("/hello", bellt.Use(
			exampleHandler,
			middlewareOne,
			middlewareTwo,
		), "GET"),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello Middleware!`))
}

func middlewareOne(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Step One")
		next.ServeHTTP(w, r)
	}
}

func middlewareTwo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Step Two")
		next.ServeHTTP(w, r)
	}
}
```

## Parameterized Routes

Route parameters must be passed using `{}` as scope limiter

```go
router.HandleFunc("/hello/{name}", handlerFunc, "GET")

router.HandleGroup("/api", 
	SubHandleFunc("/item/{id}", handlerFunc, "GET")
)
```

### Route Variables

RouteVariables used to capture and store parameters passed to built routes.

Need to pass the *Request of the HandlerFunc used in the HandleFunc method.

```go
/*
	r = *Request of the HandlerFunc
*/
rv := bellt.RouteVariables(r)
```

The declaration must be made within the HandlerFunc

```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)
	/*[...]*/
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!"))
}
```

#### GetVar

GetVar returns the parameter value of the route

```go
/*
	r = *Request of the HandlerFunc
	param = Parameter name string
*/
rv := bellt.RouteVariables(r)

rv.GetVar(param)
```

```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`Hello %v gopher!`, rv.GetVar("color")))))
}
```

The complete implementation of parameterized routes should look like this:

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {

	router := bellt.NewRouter()

	router.HandleFunc("/contact/{id}/{user}", exampleHandler, "GET")

	router.HandleGroup("/api",
		router.SubHandleFunc("/check/{id}/{user}", exampleHandler, "GET"),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"id": %v, "user": %v}`, rv.GetVar("user"), rv.GetVar("id"))))
}
```


# Full Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt"
)

func main() {

	router := bellt.NewRouter()

	router.HandleFunc("/contact/{id}/{user}", bellt.Use(
		exampleHandler,
		middlewareOne,
		middlewareTwo,
	), "GET")

	router.HandleFunc("/contact", bellt.Use(
		exampleNewHandler,
		middlewareOne,
		middlewareTwo,
	), "GET")

	router.HandleGroup("/api",
		router.SubHandleFunc("/check", bellt.Use(
			exampleNewHandler,
			middlewareOne,
			middlewareTwo,
		), "GET"),
		router.SubHandleFunc("/check/{id}/{user}", bellt.Use(
			exampleHandler,
			middlewareOne,
			middlewareTwo,
		), "GET"),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {

	rv := bellt.RouteVariables(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"id": %v, "user": %v}`, rv.GetVar("user"), rv.GetVar("id"))))
}

func exampleNewHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"msg": "Works"}`))
}

func middlewareOne(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Step One")

		next.ServeHTTP(w, r)
	}
}

func middlewareTwo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Step Two")

		next.ServeHTTP(w, r)
	}
}
```

# Benchmark

Applying requisition performance tests, the following results were obtained, showing the initial potency of the Bellt package.

<p align="center">
    <img width="800" src="./assets/benchmark.png">
</p>

# Author
Guilherme Caruso  [@guicaruso_](https://twitter.com/guicaruso_) on twitter

# Presentation
Guilherme Caruso - Cabify- GolangSP Meetup 2 - 21/03/2019 - São Paulo /Brazil  

Slides - [Construindo Rotas Parametrizadas em GO](https://www.slideshare.net/guimartinscaruso/criando-rotas-parametrizadas-em-go)

Video - [GolangSP Meetup 2](https://www.youtube.com/watch?v=nxsfyadxzmI)

# License
MIT licensed. See the LICENSE file for details.
