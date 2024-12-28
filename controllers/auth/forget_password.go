package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// ForgetPassword generates OTP for password reset
func ForgetPassword(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var dbUser models.UserRegistration
	err := utils.DB.QueryRow("SELECT username FROM users WHERE username = ?", req.Username).
		Scan(&dbUser.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate OTP
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(900000) + 100000)

	// Update OTP and set forget_is_verified to false
	_, err = utils.DB.Exec("UPDATE users SET otp = ?, forget_is_verified = ? WHERE username = ?", otp, false, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
		return
	}

	// Send OTP
	if err := utils.SendVerificationOTP(req.Username, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}
