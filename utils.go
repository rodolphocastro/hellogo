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
const minikubeUnavailableMessage = "minikube is unavailable, skipping"

var utilsLogger *zap.Logger = InitializeLogger()

// SkipTestIfMinikubeIsUnavailable Skips a test if the current environment doesn't have Minikube.
func SkipTestIfMinikubeIsUnavailable(t *testing.T) {
	isMinikubeAvailable, _ := GetMinikubeStatus()
	if !isMinikubeAvailable {
		t.Skip(minikubeUnavailableMessage)
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
	k8sLogger := utilsLogger.With(
		zap.String("k8sManifesto", pathToK8s),
	)
	k8sLogger.Info("attempting to spin up new k8s manifestos")

	minikubeIsAvailable, errStatus := GetMinikubeStatus()
	if !minikubeIsAvailable {
		k8sLogger.Info("minikube check returned unavailable", zap.String("minikubeStatus", errStatus))
		t.Skip(minikubeUnavailableMessage)
	}

	k8sLogger.Info("minikube is available")
	waitTime := time.Second
	if len(timeToWait) != 0 {
		waitTime = timeToWait[0]
	}

	pathToK8s = getPathOrDefault(pathToK8s)
	k8sLogger.Info("applying dev config")
	applyDevConfig(t)
	kubeApply := exec.Command("kubectl", "apply", "-f", pathToK8s)
	k8sLogger.Info("applying manifesto")
	err := kubeApply.Run()
	if err != nil {
		k8sLogger.Error("unexpected error while applying the manifesto", zap.Error(err))
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
	k8sLogger := utilsLogger.With(
		zap.String("k8sManifesto", pathToK8s),
	)
	k8sLogger.Info("attempting to clean up an existing k8s")

	minikubeIsAvailable, errStatus := GetMinikubeStatus()
	if !minikubeIsAvailable {
		k8sLogger.Info("minikube check returned unavailable", zap.String("minikubeStatus", errStatus))
		t.Skip(minikubeUnavailableMessage)
	}
	pathToK8s = getPathOrDefault(pathToK8s)

	k8sLogger.Info("deleting the selected manifesto")
	kubeDelete := exec.Command("kubectl", "delete", "-f", pathToK8s, "-f", pathToDevConfigs)
	err := kubeDelete.Run()
	if err != nil {
		k8sLogger.Error("an unexpected error happened while deleting the manifesto", zap.Error(err))
		t.Errorf("error while cleaning up %v -  %v", pathToK8s, err)
	}
}

// applyDevConfig Applies the ConfigMap required for development environments.
func applyDevConfig(t *testing.T) {
	minikubeIsAvailable, _ := GetMinikubeStatus()
	if !minikubeIsAvailable {
		t.Skip(minikubeUnavailableMessage)
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
	//goland:noinspection ALL
	defer logger.Sync()
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
