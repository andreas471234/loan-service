package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default configuration")
	}

	// Initialize database
	initDB()

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Setup routes
	setupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting loan service on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("loan_service.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Loan{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
}

func setupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "loan-service"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Loan routes
		loans := api.Group("/loans")
		{
			loans.GET("/", getLoans)
			loans.GET("/:id", getLoan)
			loans.POST("/", createLoan)
			loans.PUT("/:id", updateLoan)
			loans.DELETE("/:id", deleteLoan)
			loans.PUT("/:id/approve", approveLoan)
			loans.PUT("/:id/invest", investLoan)
			loans.PUT("/:id/disburse", disburseLoan)
			loans.GET("/:id/transitions", getLoanTransitions)
		}
	}
} 