package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"backend/models"
	"backend/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ResendOTP(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
	}

	// Bind the incoming JSON request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.UserRegistration
	// Check if the user exists in the database
	err := utils.DB.QueryRow("SELECT username, otp, isVerified FROM users WHERE username = ?", user.Username).
		Scan(&dbUser.Username, &dbUser.OTP, &dbUser.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Generate a new OTP
	rand.Seed(time.Now().UnixNano())
	newOTP := strconv.Itoa(rand.Intn(900000) + 100000)

	// Update the OTP in the database
	_, err = utils.DB.Exec("UPDATE users SET otp = ? WHERE username = ?", newOTP, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend OTP"})
		return
	}

	// Send the new OTP to the user (Assuming sendVerificationOTP is a predefined function)
	if err := utils.SendVerificationOTP(user.Username, newOTP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "New OTP sent successfully"})
}
