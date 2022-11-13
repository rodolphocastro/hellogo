package main

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

const defaultReply = "Hello from a test!"
const serverPort = 8762

// Using the default net/http module we can set up a http server by using funcs that implement the http.HandlerFunc
// interface, this then allows one to map a string route to a specific func that is meant to handle its request
// and writes off to the http.Response writer.
func TestServeGets(t *testing.T) {
	// Arrange
	logger := initializeZap()
	serverAddress := fmt.Sprintf(":%v", serverPort)

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
	res, err := http.Get(fmt.Sprintf("http://localhost%v/hello", serverAddress))

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
