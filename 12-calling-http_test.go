package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

const todoApiEndpoint = "https://jsonplaceholder.typicode.com/"
const todosBaseUrl = "todos/"

// Checks if a request was successful or not
func isRequestSuccessful(r *http.Response) bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// Reads a response body as a string
func readResponseBodyAsString(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	return string(body), err
}

// By using the http module we're able to invoke all the HTTP Methods upon an API.
// https://gobyexample.com/http-clients
func TestCallGet(t *testing.T) {
	allTodos := fmt.Sprintf("%v%v", todoApiEndpoint, todosBaseUrl)
	// Hitting all the todos on the placeholder api
	response, err := http.Get(allTodos)
	if err != nil {
		t.Errorf("Expected no errors getting %v but got %v", allTodos, err)
	}

	// deferring to allow the request to complete
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Error("Expected no errors while defer the request, but an error happened")
		}
	}()

	// Checking the response's status code
	status := response.StatusCode
	if !isRequestSuccessful(response) {
		t.Errorf("Expected status to be Ok but found %v", status)
	}

	// getting the response's body and its contents
	todos, err := readResponseBodyAsString(response)
	if err != nil {
		t.Error("Expected no errors while parsing the response, but an error happened")
	}

	if todos == "" {
		t.Error("Expected response to not be blank but got blank")
	}
}
