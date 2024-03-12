package config

import (
	"os"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	// backup the original environment variables
	envMap := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envMap[pair[0]] = pair[1]
	}

	// set the environment variables for testing
	os.Setenv("APP_SECRET", "test_secret")
	os.Setenv("PORT", "5001")

	// load the test configuration
	appConfig, err := LoadConfig("../config.yaml")
	if err != nil {
		t.Fatalf("Failed to load test configuration: %v", err)
	}

	expectedAppSecret := "test_secret"
	if appConfig.AppSecret != expectedAppSecret {
		t.Errorf("Expected AppSecret to be '%s', got '%s'", expectedAppSecret, appConfig.AppSecret)
	}

	expectedPort := 5001
	if appConfig.Port != expectedPort {
		t.Errorf("Expected Port to be '%d', got '%d'", expectedPort, appConfig.Port)
	}

	// recover the original environment variables
	for k, v := range envMap {
		os.Setenv(k, v)
	}
}
