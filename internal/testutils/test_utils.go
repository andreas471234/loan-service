package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"loan-service/internal/domain"
	"loan-service/internal/handler"
	"loan-service/internal/middleware"
	"loan-service/internal/repository"
	"loan-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSetup contains test configuration
type TestSetup struct {
	Router *gin.Engine
	DB     *gorm.DB
	Server *httptest.Server
}

// SetupTestDB creates a test database
func SetupTestDB() *gorm.DB {
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

// SetupTestRouter creates a test router without setting up routes
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// SetupTestServer creates a test server with in-memory database (for integration tests)
func SetupTestServer() *TestSetup {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test database
	testDB := SetupTestDB()

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "loan-service"})
	})

	// Initialize dependencies
	loanRepo := repository.NewLoanRepository(testDB)
	loanService := service.NewLoanService(loanRepo)
	loanHandler := handler.NewLoanHandler(loanService)

	// API routes
	api := router.Group("/api/v1")
	{
		// Loan routes
		loans := api.Group("/loans")
		{
			loans.GET("/", loanHandler.GetLoans)
			loans.GET("/:id", loanHandler.GetLoan)
			loans.POST("/", loanHandler.CreateLoan)
			loans.PUT("/:id", loanHandler.UpdateLoan)
			loans.DELETE("/:id", loanHandler.DeleteLoan)
			loans.PUT("/:id/approve", loanHandler.ApproveLoan)
			loans.PUT("/:id/invest", loanHandler.InvestLoan)
			loans.PUT("/:id/disburse", loanHandler.DisburseLoan)
			loans.GET("/:id/transitions", loanHandler.GetLoanTransitions)
		}
	}

	// Create test server
	server := httptest.NewServer(router)

	return &TestSetup{
		Router: router,
		DB:     testDB,
		Server: server,
	}
}

// MakeRequest is a helper function to make HTTP requests
func MakeRequest(method, url string, body interface{}) (*http.Response, error) {
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
func Float64Ptr(v float64) *float64 {
	return &v
}

func StringPtr(v string) *string {
	return &v
}
