package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const defaultPathToK8s = "./k8s.yaml"
const pathToDevConfigs = "./environments/development/config.yml"

const ciEnvKey = "CI"

// SkipTestIfCI Skips a test if the current environment is a CI pipeline.
func SkipTestIfCI(t *testing.T) {
	if isEnvironmentCI() {
		t.Skip("Skipping this test - we're running in a CI environment")
	}
}

// getMinikubeIp gets the Minikube IP from the OS' console. If minikube is unavailable it'll return an empty string.
func getMinikubeIp() string {
	command := exec.Command("minikube", "ip")
	byteResult, err := command.Output()
	if err != nil {
		return ""
	}
	result := string(byteResult)
	return strings.TrimSpace(result)
}

// isEnvironmentCI checks if the current environment is a Continuous Integration pipeline.
func isEnvironmentCI() bool {
	return os.Getenv(ciEnvKey) != ""
}

// SpinUpK8s Quick and Dirty way to spin up the deployment - invoking kubectl in the os' console.
func SpinUpK8s(t *testing.T, pathToK8s string) {
	pathToK8s = getPathOrDefault(pathToK8s)
	applyDevConfig(t)
	kubeApply := exec.Command("kubectl", "apply", "-f", pathToK8s)
	if kubeApply.Run() != nil {
		t.Error("Error while spinning up MongoDb")
	}
}

// getPathOrDefault returns the current path or a default one in case none is set.
func getPathOrDefault(pathToK8s string) string {
	if pathToK8s == "" {
		pathToK8s = defaultPathToK8s
	}
	return pathToK8s
}

// CleanUpK8s Quick and Dirty way to delete the deployment - invoking kubectl in the os' console.
func CleanUpK8s(t *testing.T, pathToK8s string) {
	pathToK8s = getPathOrDefault(pathToK8s)

	kubeDelete := exec.Command("kubectl", "delete", "-f", pathToK8s, "-f", pathToDevConfigs)
	if kubeDelete.Run() != nil {
		t.Error("Error while cleaning up MongoDb")
	}
}

// applyDevConfig Applies the ConfigMap required for development environments.
func applyDevConfig(t *testing.T) {
	kubeConfig := exec.Command("kubectl", "apply", "-f", pathToDevConfigs)
	if kubeConfig.Run() != nil {
		t.Error("Error while applying Dev's ConfigMaps")
	}
}
