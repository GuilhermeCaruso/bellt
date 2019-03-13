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