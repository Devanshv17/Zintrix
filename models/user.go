package models

import "time"

// User struct represents the schema for the "users" table
type UserRegistration struct {
	ID               int       `json:"id"`
	Username         string    `json:"username" binding:"required"`
	Password         string    `json:"password" binding:"required"`
	IsVerified       bool      `json:"is_verified"`
	ForgetIsVerified bool      `json:"forget_is_verified"`
	OTP              string    `json:"otp"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
