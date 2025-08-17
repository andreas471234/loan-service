package database

import (
	"testing"

	"loan-service/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	// Test SQLite connection
	cfg := config.DatabaseConfig{
		Driver: "sqlite",
		Name:   ":memory:",
	}

	db, err := NewConnection(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Test unsupported driver
	cfg.Driver = "unsupported"
	_, err = NewConnection(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database driver")
}

func TestCloseConnection(t *testing.T) {
	// Test with nil database (should not panic)
	CloseConnection(nil)

	// Test with valid database
	cfg := config.DatabaseConfig{
		Driver: "sqlite",
		Name:   ":memory:",
	}

	db, err := NewConnection(cfg)
	require.NoError(t, err)

	CloseConnection(db)
}

func TestGetDB(t *testing.T) {
	// Test GetDB function
	_ = GetDB()
	// This might be nil if no connection was established, which is fine for testing
	// We're just testing that the function doesn't panic
}
