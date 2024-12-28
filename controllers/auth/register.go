package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(c *gin.Context) {
	var user models.UserRegistration
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	var existingUser models.UserRegistration
	err := utils.DB.QueryRow("SELECT username FROM users WHERE username = ?", user.Username).
		Scan(&existingUser.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already registered"})
		return
	}

	// Generate OTP
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(900000) + 100000)

	// Send OTP
	if err := utils.SendVerificationOTP(user.Username, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save the user to the database
	_, err = utils.DB.Exec("INSERT INTO users (username, password, isVerified, otp) VALUES (?, ?, ?, ?)",
		user.Username, string(hashedPassword), false, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully. Please verify your email to activate your account"})
}
