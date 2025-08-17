package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	v1 "loan-service/api/v1"
	"loan-service/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSetup contains test configuration
type TestSetup struct {
	router *gin.Engine
	db     *gorm.DB
	server *httptest.Server
}

// setupTestDB creates a test database
func setupTestDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto migrate the schema
	err = database.AutoMigrate(&domain.Loan{}, &domain.Investment{})
	if err != nil {
		panic("failed to migrate test database")
	}

	return database
}

// setupTestServer creates a test server with in-memory database
func setupTestServer() *TestSetup {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test database
	testDB := setupTestDB()

	// Create router
	router := gin.New()

	// Setup test routes using the new structure
	v1.SetupRoutes(router, testDB)

	// Create test server
	server := httptest.NewServer(router)

	return &TestSetup{
		router: router,
		db:     testDB,
		server: server,
	}
}

// makeRequest is a helper function to make HTTP requests
func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

// Helper functions
func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}
