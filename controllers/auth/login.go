package controllers

import (
	"log"
	"net/http"

	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login
func Login(c *gin.Context) {
	var user models.UserRegistration
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.UserRegistration
	// Fetch user data from the database
	err := utils.DB.QueryRow("SELECT username, password, isVerified FROM users WHERE username = ?", user.Username).
		Scan(&dbUser.Username, &dbUser.Password, &dbUser.IsVerified)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !dbUser.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified"})
		return
	}
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		// Log if the password comparison fails
		log.Printf("Login: Password comparison failed for user %s. Error: %v", user.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
