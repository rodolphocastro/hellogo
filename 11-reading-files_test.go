package main

import (
	"io"
	"os"
	"testing"
)

const sampleFileName = "./11-files-sample.json"
const newFileName = "./11-files-new-file.json"
const newFileContent = "I like turtles"

// Reading a file's content with the os module
// https://gobyexample.com/reading-files
func TestReadFromOs(t *testing.T) {
	// to read everything from a file we need to use the os module
	contents, err := os.ReadFile(sampleFileName)
	if err != nil {
		t.Error("Something went wrong when reading from the sample")
	}
	contentAsString := string(contents)

	if contentAsString == "" {
		t.Errorf("Expected content to not be empty but found %v", contentAsString)
	}
}

// Reading a file's content with both os and io modules
// https://gobyexample.com/reading-files
func TestReadAndParse(t *testing.T) {
	// getting a Reader that points to a file (for more granular operations on files)
	file, err := os.Open(sampleFileName)
	if err != nil {
		t.Error("Something went wrong when reading from the sample")
	}

	// Reading all the contents from the file's Reader
	contents, err := io.ReadAll(file)
	if err != nil {
		t.Error("Something went wrong when reading from the sample")
	}

	if string(contents) == "" {
		t.Error("Expected content to not be empty but found empty")
	}
}

// Deleting, then creating and writing to a file
// https://gobyexample.com/writing-files
func TestDeleteAndWrite(t *testing.T) {
	// Attempting to delete a file
	_ = os.Remove(newFileName) // Not error checking because we don't care if the deletion worked

	// Creating a new file by dumping all the bytes from a string
	err := os.WriteFile(newFileName, []byte(newFileContent), 0644)
	if err != nil {
		t.Error("An error happened while creating the file")
	}

	// Reading the contents back for sanity checking
	contents, err := os.ReadFile(newFileName)
	if err != nil {
		t.Error("An error happened while reading from the file")
	}

	if string(contents) != newFileContent {
		t.Errorf("Expected %v but read %v from the file", newFileContent, contents)
	}

	err = os.Remove(newFileName)
	if err != nil {
		t.Error("An error happened while deleting the file")
	}
}
