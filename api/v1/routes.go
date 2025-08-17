package v1

import (
	"loan-service/internal/handler"
	"loan-service/internal/middleware"
	"loan-service/internal/repository"
	"loan-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "loan-service"})
	})

	// Initialize dependencies
	loanRepo := repository.NewLoanRepository(db)
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
}
