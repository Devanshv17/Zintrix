package controllers

import (
	"net/http"

	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// VerifyOTP handles OTP verification
func VerifyOTP(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		OTP      string `json:"otp" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if OTP matches
	var dbUser models.UserRegistration
	err := utils.DB.QueryRow("SELECT otp FROM users WHERE username = ?", req.Username).
		Scan(&dbUser.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or OTP"})
		return
	}

	if dbUser.OTP != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Update isVerified field
	_, err = utils.DB.Exec("UPDATE users SET isVerified = ? WHERE username = ?", true, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}
