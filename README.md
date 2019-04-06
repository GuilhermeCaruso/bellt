# Bellt  
> Simple Golang HTTP router

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/GuilhermeCaruso/bellt)](https://goreportcard.com/report/github.com/GuilhermeCaruso/bellt) [![codecov](https://codecov.io/gh/GuilhermeCaruso/bellt/branch/master/graph/badge.svg)](https://codecov.io/gh/GuilhermeCaruso/bellt) [![Build Status](https://travis-ci.com/GuilhermeCaruso/bellt.svg?branch=master)](https://travis-ci.com/GuilhermeCaruso/bellt) ![GitHub](https://img.shields.io/badge/golang%20->=1.7-blue.svg) [![GoDoc](https://godoc.org/github.com/GuilhermeCaruso/bellt?status.svg)](https://godoc.org/github.com/GuilhermeCaruso/bellt) 

<p align="left">
    <img width="150" src="./assets/logo.png">
</p>

Bellt Package implements a request router with the aim of managing controller actions based on fixed and parameterized routes.

The project so far has the following functionalities:

* Standard definition of route "/ health", in order to prepare the service developed with bellt to act as microservice.
* Providing the creation of parameterized routes, simple or segmented (groups).
* All requests can be made through fixed patterns, querystrings and parameters.
* Obtaining the requisition parameters in the controller functions.
---
# Summary
 * [Install](#install)
 * [Guide](#guide)
	* [Router](#router)
	* [Methods](#router-methods)
		* [HandleFunc](#handleFunc)
		* [HandleGroup](#handleFunc)
 * Parameterized Routes
 * Middleware
 * Full Example

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

## Router Methods

#### HandleFunc   

HandleFunc function responsible for initializing a common route or built through the Router. All non-grouped routes must be initialized by this method.

```go
    /*
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
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
- HandleFunc with Middlewares
```go
    /*
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
    */

    router.HandleFunc(path, bellt.Use(
        handlerFunc,
        middlewareOne,
        middlewareTwo,
    ), methods)
   
```

- HandleGroup && SubHandleFunc    

HandleGroup used to create and define a group of sub-routes.

SubHandleFunc is responsible for initializing a common or built route. Its use must be made within the scope of the HandleGroup() method, where the main path will be declared.
```go

    /*
        mainPath - String route grouper
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
    */

    router.HandleGroup(mainPath,
        router.SubHandleFunc(path, handlerFunc, methods),
        router.SubHandleFunc(path, bellt.Use(
            handlerFunc,
            middlewareOne,
            middlewareTwo,
        ), methods),
    )
    
```

## Route Params

To use parameters in routes we must group the variables in braces.
```go
    router.HandleFunc("/user/{id}", handlerFunc, methods)
    router.HandleGroup("/api",
        router.SubHandleFunc("/product/{id}", handlerFunc, methods),
        router.SubHandleFunc("/product/{id}/{categorie}", bellt.Use(
            handlerFunc,
            middlewareOne,
            middlewareTwo,
        ), methods),
    )
```
## Get Route Params
```go
    router.HandleFunc("/user/{id}/{user}", exampleFunction, "GET")

    func exampleFunction(w http.ResponseWriter, r *http.Request) {
        rv := bellt.RouteVariables(r)

        fmt.Println(rv.GetVar("user"))
        fmt.Println(rv.GetVar("id"))
    }
    // Guilherme
    // 123

```
    

# Examples

Let's start our simple router application.

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

	router.HandleFunc("/health", healthApplication , "GET", "PUT")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

## Full Example
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
	), "GET", "PUT")

	router.HandleFunc("/contact", bellt.Use(
		exampleNewHandler,
		middlewareOne,
		middlewareTwo,
	), "GET", "PUT")

	router.HandleGroup("/api",
		router.SubHandleFunc("/check", bellt.Use(
			exampleNewHandler,
			middlewareOne,
			middlewareTwo,
		), "GET", "PUT"),
		router.SubHandleFunc("/check/{id}/{user}", bellt.Use(
			exampleHandler,
			middlewareOne,
			middlewareTwo,
		), "GET", "PUT"),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"id": %v, "user": %v}`, rv.GetVar("user"), rv.GetVar("id"))))
}

func exampleNewHandler(w http.ResponseWriter, r *http.Request) {
	rv := bellt.RouteVariables(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
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

# Author
Guilherme Caruso  [@guicaruso_](https://twitter.com/guicaruso_) on twitter

# Presentation
Guilherme Caruso - Cabify- GolangSP Meetup 2 - 21/03/2019 - São Paulo /Brazil  

Slides - [Construindo Rotas Parametrizadas em GO](https://www.slideshare.net/guimartinscaruso/criando-rotas-parametrizadas-em-go)

Video - [GolangSP Meetup 2](https://www.youtube.com/watch?v=nxsfyadxzmI)
# License
MIT licensed. See the LICENSE file for details.
