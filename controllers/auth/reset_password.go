package controllers

import (
	"net/http"

	"backend/models"
	"backend/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	// Bind the incoming JSON request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.UserRegistration
	// Check if the username exists and the OTP has been verified
	err := utils.DB.QueryRow("SELECT username, forget_is_verified FROM users WHERE username = ?", req.Username).
		Scan(&dbUser.Username, &dbUser.ForgetIsVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if the OTP has been verified
	if !dbUser.ForgetIsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your OTP first"})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update the password and set forget_is_verified to false
	_, err = utils.DB.Exec("UPDATE users SET password = ?, forget_is_verified = ? WHERE username = ?",
		string(hashedPassword), false, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
