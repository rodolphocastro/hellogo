package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const defaultReply = "Hello from a test!"
const serverPort = 8762

// gets the server's hosting address
func getServerRealAddress(serverAddress string) string {
	return fmt.Sprintf("http://localhost%v/hello", serverAddress)
}

// gets the server's binding address
func getServerBindingAddress() string {
	serverAddress := fmt.Sprintf(":%v", serverPort)
	return serverAddress
}

// Using the default net/http module we can set up a http server by using funcs that implement the http.HandlerFunc
// interface, this then allows one to map a string route to a specific func that is meant to handle its request
// and writes off to the http.Response writer.
func TestServeGets(t *testing.T) {
	// Arrange
	logger := initializeZap()
	serverAddress := getServerBindingAddress()

	// spinning a new goRoutine to serve the server (pun intended)
	go func() {
		logger.Info("spinning up a new goRoutine for the HttpServer")
		http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
			// On goLang's net/http module all we need to do is implement the interface for http.Handler
			defer logger.Info("done responding to the message")
			logger.Info("received a new request",
				zap.String("host", request.Host),
			)

			_, err := fmt.Fprintf(writer, defaultReply)
			if err != nil {
				logger.Error("an error happened while replying", zap.Error(err))
			}
		})

		err := http.ListenAndServe(serverAddress, nil)
		if err != nil {
			logger.Error("something didn't go as expected", zap.Error(err))
		}
	}()

	// Act
	res, err := http.Get(getServerRealAddress(serverAddress))

	// Assert
	if err != nil {
		t.Errorf("expected no errors but got %v", err)
	}

	if res.StatusCode != 200 {
		t.Error("expected an Ok response but got something else")
	}

	bodyContents, _ := io.ReadAll(res.Body)
	stringResult := string(bodyContents)
	if stringResult != defaultReply {
		t.Errorf("expected %v but found %v", defaultReply, stringResult)
	}
}

// Using net/http/httptest we can easily create mocks and stubs to test the most common scenarios a Http Server
// and its handlers need to deal with.
func TestUsingHttpTestForTesting(t *testing.T) {
	// Arrange
	logger := initializeZap()
	testRequest := httptest.NewRequest(http.MethodGet, "http://example.io/something", nil)
	recorder := httptest.NewRecorder()

	getSomethingHandler := func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handler started")
		defer logger.Info("handler finished")

		_, err := fmt.Fprintf(w, defaultReply)
		if err != nil {
			logger.Error("unexpected error replying to a request", zap.Error(err))
		}
	}

	// Act
	// firing off a dummy request to the Handler
	getSomethingHandler(recorder, testRequest)

	// Assert
	gotStatusCode := recorder.Result().StatusCode
	if gotStatusCode != http.StatusOK {
		t.Errorf("expected an Ok response but got %v instead", gotStatusCode)
	}

	gotBytes, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		t.Errorf("expected no errors but found %v", err)
	}

	gotBody := string(gotBytes)
	if gotBody != defaultReply {
		t.Fatalf("expected response's body to be %v but found %v instead", defaultReply, gotBody)
	}
}

// We can access a request's context to have more information about the caller and the route itself. In this test we
// inject our own context to mock it successfully.
func TestARequestContextCanBeAccessedForMoreInformationAboutTheInvoker(t *testing.T) {
	// Arrange
	const expectedKey = "dummy"
	logger := initializeZap()
	testRequest := httptest.NewRequest(http.MethodGet, "http://example.io/something", strings.NewReader(defaultReply))
	// adding a dummy value to the request's context
	requestCtx := testRequest.Context()
	requestCtx = context.WithValue(requestCtx, expectedKey, defaultReply)
	testRequest = testRequest.WithContext(requestCtx)
	// wrapping up - now the context has the expected value!
	recorder := httptest.NewRecorder()
	results := make(chan string, 1)

	subjectHandler := func(writer http.ResponseWriter, request *http.Request) {
		logger.Info("handling an incoming request!")
		defer logger.Info("request handled, wrapping up")
		defer close(results)
		ctx := request.Context() // Since we're on a test this context is nil (empty) so there's not much we can read from
		writer.WriteHeader(http.StatusGone)
		_, err := writer.Write([]byte("oh no no no"))
		if err != nil {
			logger.Error("an unexpected error happened while writing a reply", zap.Error(err))
		}
		if err != nil {
			logger.Error("an unexpected error happened retrieving the body's contents", zap.Error(err))
		}

		contextValue := ctx.Value(expectedKey).(string)
		logger.Info("here's the context for this request", zap.Any("requestContext", ctx))
		logger.Info("publishing the default value within the context", zap.String("dummyValue", contextValue))
		results <- contextValue
	}

	// Act
	subjectHandler(recorder, testRequest)

	// Assert
	got := <-results
	if got != defaultReply {
		t.Errorf("expected %v but got %v", defaultReply, got)
	}
}
