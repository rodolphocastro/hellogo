package main

import (
	"os"
	"os/exec"
	"testing"
)

const defaultPathToK8s = "./k8s.yaml"
const ciEnvKey = "CI"

// SkipTestIfCI Skips a test if the current environment is a CI pipeline.
func SkipTestIfCI(t *testing.T) {
	if os.Getenv(ciEnvKey) != "" {
		t.Skip("Skipping this test - we're running in a CI environment")
	}
}

// SpinUpMongoK8s Quick and Dirty way to spin up the deployment - invoking kubectl in the os' console.
func SpinUpMongoK8s(t *testing.T, pathToK8s string) {
	if pathToK8s == "" {
		pathToK8s = defaultPathToK8s
	}

	kubeApply := exec.Command("kubectl", "apply", "-f", pathToK8s)
	if kubeApply.Run() != nil {
		t.Error("Error while spinning up MongoDb")
	}
}

// CleanUpMongoK8s Quick and Dirty way to delete the deployment - invoking kubectl in the os' console.
func CleanUpMongoK8s(t *testing.T, pathToK8s string) {
	if pathToK8s == "" {
		pathToK8s = defaultPathToK8s
	}
	
	kubeDelete := exec.Command("kubectl", "delete", "-f", pathToK8s)
	if kubeDelete.Run() != nil {
		t.Error("Error while cleaning up MongoDb")
	}
}
