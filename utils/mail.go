package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendVerificationOTP sends a verification OTP to the provided email address
func SendVerificationOTP(email, otp string) error {
	// Replace with your SMTP server details
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := os.Getenv("SMTP_USERNAME")    // Use environment variable for security
	senderPassword := os.Getenv("SMTP_PASSWORD") // Use environment variable for security

	// Email body
	subject := "Subject: Your Verification Code\n"
	body := fmt.Sprintf("Your OTP is: %s\n", otp)
	message := subject + "\n" + body

	// Set up authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
