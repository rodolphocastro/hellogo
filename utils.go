package main

import (
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const defaultPathToK8s = "./k8s.yaml"
const pathToDevConfigs = "./environments/development/config.yml"
const cicdPipelineEnvKey = "CI"

// SkipTestIfMinikubeIsUnavailable Skips a test if the current environment doesn't have Minikube.
func SkipTestIfMinikubeIsUnavailable(t *testing.T) {
	isMinikubeAvailable, _ := GetMinikubeStatus()
	if !isMinikubeAvailable {
		t.Skip("skipping test because minikube is unavailable")
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
	isCiCdEnvSet := os.Getenv(cicdPipelineEnvKey) != ""
	return isCiCdEnvSet
}

// SpinUpK8s Quick and Dirty way to spin up the deployment - invoking kubectl in the os' console.
func SpinUpK8s(t *testing.T, pathToK8s string, timeToWait ...time.Duration) {
	minikubeIsAvailable, _ := GetMinikubeStatus()
	if !minikubeIsAvailable {
		t.Skip("minikube isn't available, skipping")
	}

	waitTime := time.Second
	if len(timeToWait) != 0 {
		waitTime = timeToWait[0]
	}

	pathToK8s = getPathOrDefault(pathToK8s)
	applyDevConfig(t)
	kubeApply := exec.Command("kubectl", "apply", "-f", pathToK8s)
	if kubeApply.Run() != nil {
		t.Errorf("error while spinning up environment %v", pathToK8s)
	}
	time.Sleep(waitTime)
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
	minikubeIsAvailable, _ := GetMinikubeStatus()
	if !minikubeIsAvailable {
		t.Skip("minikube isn't available, skipping")
	}
	pathToK8s = getPathOrDefault(pathToK8s)

	kubeDelete := exec.Command("kubectl", "delete", "-f", pathToK8s, "-f", pathToDevConfigs)
	if kubeDelete.Run() != nil {
		t.Error("Error while cleaning up MongoDb")
	}
}

// applyDevConfig Applies the ConfigMap required for development environments.
func applyDevConfig(t *testing.T) {
	minikubeIsAvailable, _ := GetMinikubeStatus()
	if !minikubeIsAvailable {
		t.Skip("minikube isn't available, skipping")
	}
	kubeConfig := exec.Command("kubectl", "apply", "-f", pathToDevConfigs)
	if kubeConfig.Run() != nil {
		t.Error("Error while applying Dev's ConfigMaps")
	}
}

// InitializeLogger initializes a Zap logger and returns it based on the environment.
// default behavior is a Production logger for CI and Development logger for non-CI environments.
func InitializeLogger() *zap.Logger {
	var logger *zap.Logger
	var err error

	if isEnvironmentCI() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Error while setting up zap")
	}

	return logger
}

// GetMinikubeStatus gets the current status of the Minikube cluster (true if running, false otherwise)
// and its details status
func GetMinikubeStatus() (bool, string) {
	command := exec.Command("minikube", "status")
	result, err := command.Output()
	minikubeDetailsStatus := string(result)
	isRunning := err == nil
	return isRunning, minikubeDetailsStatus
}
