package bellt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Error("Router initialize with errors!")
	}

}

func TestGetRouter(t *testing.T) {
	router := NewRouter()
	sameRouter := getRouter()
	if router != sameRouter {
		t.Error("Router with colapses!")
	}
}

func TestBuiltRoute(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/user/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		rv := RouteVariables(r)
		id := fmt.Sprintf("%v", rv.GetVar("id"))
		w.Write([]byte(id))
	}, "GET")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectBuiltRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `123`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRoute(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("works"))
	}, "GET")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("works"))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `works`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func oneMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("one")

		next.ServeHTTP(w, r)
	}
}

func TestMiddleware(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/user/789", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/user/{id}", Use(
		func(w http.ResponseWriter, r *http.Request) {
			rv := RouteVariables(r)
			id := fmt.Sprintf("%v", rv.GetVar("id"))
			w.Write([]byte(id))
		},
		oneMiddleware,
	), "GET")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectBuiltRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `789`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGroupRoute(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/api/user/789", nil)
	if err != nil {
		t.Fatal(err)
	}
	router.HandleGroup("api/",
		router.SubHandleFunc("user/{id}", func(w http.ResponseWriter, r *http.Request) {
			rv := RouteVariables(r)
			id := fmt.Sprintf("%v", rv.GetVar("id"))
			w.Write([]byte(id))
		}, "GET"),
	)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectBuiltRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `789`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/user/456", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		rv := RouteVariables(r)
		id := fmt.Sprintf("%v", rv.GetVar("id"))
		w.Write([]byte(id))
	}, "POST")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectBuiltRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"error": "The method for this route doesnt exist"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthModel(t *testing.T) {

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthApplication)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": "Server running"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
