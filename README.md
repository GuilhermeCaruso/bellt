<p align="center">
    <img width="200" src="./logo.png">
</p>


# Bellt  
> Simple Golang HTTP router

## Setup

To get Bellt

##### > Go CLI
```sh
go get github.com/GuilhermeCaruso/bellt
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

## Router Methods

- HandleFunc   

HandleFunc function responsible for initializing a common route or built through the Router. All non-grouped routes must be initialized by this method.

```go
    router.HandleFunc(path, handlerFunc, methods)
    /*
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
    */
```
- HandleFunc with Middlewares
```go
    router.HandleFunc(path, bellt.Use(
        handlerFunc,
        middlewareOne,
        middlewareTwo,
    ), methods)
    /*
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
    */
```

- HandleGroup && SubHandleFunc    

HandleGroup used to create and define a group of sub-routes.

SubHandleFunc is responsible for initializing a common or built route. Its use must be made within the scope of the HandleGroup() method, where the main path will be declared.
```go
    router.HandleGroup(mainPath,
        router.SubHandleFunc(path, handlerFunc, methods),
        router.SubHandleFunc(path, bellt.Use(
            handlerFunc,
            middlewareOne,
            middlewareTwo,
        ), methods),
    )
    /*
        mainPath - String route grouper
        path - Endpoint string
        handlerFunc - function that will be called on the request
        methods - Slice for endpoint methods ("GET", "POST", "PUT", "DELETE")
    */
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

	router.HandleFunc("/healt", healthApplication , "GET", "PUT")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```


# Author
Guilherme Caruso  [@guicaruso_](https://twitter.com/guicaruso_) on twitter

# License
BSD licensed. See the LICENSE file for details.