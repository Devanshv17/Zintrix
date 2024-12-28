package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/controllers"
	auth "backend/controllers/auth"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	// Initialize the database connection
	utils.InitializeDatabase(dsn)
}

func main() {
	defer controllers.SendShutdownNotification()
	// Initialize the Gin router
	router := gin.Default()

	// Middleware for logging and recovering from panics
	router.Use(gin.Recovery())

	// Public routes
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.POST("/verify-otp", auth.VerifyOTP)
	router.POST("/resend-otp", auth.ResendOTP)
	router.POST("/forgot-password", auth.ForgetPassword)
	router.POST("/verify-reset", auth.VerifyOTPForPasswordReset)
	router.POST("/reset-password", auth.ResetPassword)

	// Create the server and bind to port
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
