package router_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/trencat/goutils/http/router"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Test handler\n")
}

func middleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Middleware acting before handler\n")
		fn(w, r)
		fmt.Fprint(w, "Middleware acting after handler\n")
	}
}

func funnyMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "FunnyMiddleware acting before handler\n")
		fn(w, r)
		fmt.Fprint(w, "FunnyMiddleware acting after handler\n")
	}
}

// Init a server at given address
func server(addr string, mux *http.ServeMux) {
	// Define new router
	r := make(router.Router)

	// Request with no middleware
	r.HandleFunc("/nomiddleware", handler)

	// Request with middleware
	r.Add("/middleware", middleware)
	r.Add("/middleware", funnyMiddleware)
	r.HandleFunc("/middleware", handler)

	r.Build(mux)

	// Run server in a new thread
	if mux == nil {
		go http.ListenAndServe(addr, nil)
	} else {
		go http.ListenAndServe(addr, mux)
	}

	// Wait a second to make sure is up
	time.Sleep(time.Duration(1) * time.Second)
}

func ExampleRouter() {
	// Run a server to listen to HTTP requests
	server(":8080", nil)

	// Do request. Omitting error handling for clarity
	response, _ := http.Get("http://127.0.0.1:8080/nomiddleware")
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s", body)

	// Do request. Omitting error handling for clarity
	response, _ = http.Get("http://127.0.0.1:8080/middleware")
	body, _ = ioutil.ReadAll(response.Body)
	fmt.Printf("%s", body)
}

// GetRequest performs a get request and return response body
func GetRequest(addr string, t *testing.T) []byte {
	t.Helper()

	response, err := http.Get(addr)
	if err != nil {
		t.Errorf("Error in GetRequest. Got \"%+v\". Expected nil error", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error in GetRequest. Got %+v. Expected nil error", err)
	}

	return body
}

var serverPort string
var portDefault string
var portMux string

func TestMain(m *testing.M) {
	portDefault, portMux = ":8081", ":8082"

	// Run server
	server(portDefault, nil)

	// Run server in a custom multiplexer
	mux := http.NewServeMux()
	server(portMux, mux)

	os.Exit(m.Run())
}

func TestRouterNoMiddleware(t *testing.T) {
	requestAddr := fmt.Sprintf("http://127.0.0.1%s/nomiddleware", portDefault)
	response := GetRequest(requestAddr, t)

	resp := fmt.Sprintf("%s", response)
	expected := "Test handler\n"
	if resp != expected {
		t.Errorf("Got %s, expected \"%s\"", resp, expected)
	}
}

func TestRouterMiddleware(t *testing.T) {
	requestAddr := fmt.Sprintf("http://127.0.0.1%s/middleware", portDefault)
	response := GetRequest(requestAddr, t)

	resp := fmt.Sprintf("%s", response)
	expected := "Middleware acting before handler\n"
	expected += "FunnyMiddleware acting before handler\nTest handler\n"
	expected += "FunnyMiddleware acting after handler\nMiddleware acting after handler\n"
	if resp != expected {
		t.Errorf("Got %s, expected \"%s\"", resp, expected)
	}
}

func TestRouterMuxNoMiddleware(t *testing.T) {
	requestAddr := fmt.Sprintf("http://127.0.0.1%s/nomiddleware", portMux)
	response := GetRequest(requestAddr, t)

	resp := fmt.Sprintf("%s", response)
	expected := "Test handler\n"
	if resp != expected {
		t.Errorf("Got %s, expected \"%s\"", resp, expected)
	}
}

func TestRouterMuxMiddleware(t *testing.T) {
	requestAddr := fmt.Sprintf("http://127.0.0.1%s/middleware", portMux)
	response := GetRequest(requestAddr, t)

	resp := fmt.Sprintf("%s", response)
	expected := "Middleware acting before handler\n"
	expected += "FunnyMiddleware acting before handler\nTest handler\n"
	expected += "FunnyMiddleware acting after handler\nMiddleware acting after handler\n"
	if resp != expected {
		t.Errorf("Got %s, expected \"%s\"", resp, expected)
	}
}

//TODO: Check that nil code is returned on router.Build()
//TODO: Add example of router with no handlerfunc but yes middleware raising an error
