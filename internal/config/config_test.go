package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Test with default values
	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "sqlite", config.Database.Driver)
	assert.Equal(t, "loan_service.db", config.Database.Name)
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("PORT", "9090")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("SERVER_READ_TIMEOUT", "30")
	os.Setenv("SERVER_WRITE_TIMEOUT", "30")
	os.Setenv("SERVER_IDLE_TIMEOUT", "300")

	// Clean up after test
	defer func() {
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("PORT")
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("SERVER_READ_TIMEOUT")
		os.Unsetenv("SERVER_WRITE_TIMEOUT")
		os.Unsetenv("SERVER_IDLE_TIMEOUT")
	}()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "production", config.Environment)
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "postgres", config.Database.Driver)
	assert.Equal(t, "test_db", config.Database.Name)
	assert.Equal(t, 30, int(config.Server.ReadTimeout.Seconds()))
	assert.Equal(t, 30, int(config.Server.WriteTimeout.Seconds()))
	assert.Equal(t, 300, int(config.Server.IdleTimeout.Seconds()))
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default")
	assert.Equal(t, "test_value", value)

	// Test with non-existing environment variable
	value = getEnv("NON_EXISTENT_VAR", "default_value")
	assert.Equal(t, "default_value", value)
}
