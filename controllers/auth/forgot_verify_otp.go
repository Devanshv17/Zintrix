package controllers

import (
	"net/http"

	"backend/models"
	"backend/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func VerifyOTPForPasswordReset(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		OTP      string `json:"otp" binding:"required"`
	}

	// Bind the incoming JSON request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.UserRegistration
	// Check if the OTP matches
	err := utils.DB.QueryRow("SELECT username, otp, forget_is_verified FROM users WHERE username = ?", req.Username).
		Scan(&dbUser.Username, &dbUser.OTP, &dbUser.ForgetIsVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Verify if the OTP matches and is not expired (you can implement expiration logic here)
	if dbUser.OTP != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Update the user's forget_is_verified status to true after OTP verification
	_, err = utils.DB.Exec("UPDATE users SET forget_is_verified = 1 WHERE username = ?", req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified. You can now reset your password"})
}
