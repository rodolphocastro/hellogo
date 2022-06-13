package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

const todoApiEndpoint = "https://jsonplaceholder.typicode.com/"
const todosBaseUrl = "todos/"

// By using the http module we're able to invoke all the HTTP Methods upon an API.
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
	if status != http.StatusOK {
		t.Errorf("Expected status to be Ok but found %v", status)
	}

	// getting the response's body and its contents
	todos, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("Expected no errors while parsing the response, but an error happened")
	}

	if string(todos) == "" {
		t.Error("Expected response to not be blank but got blank")
	}
}
