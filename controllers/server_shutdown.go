package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendShutdownNotification() {
	// SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	from := smtpUsername         // Change this to your email address
	to := "devanshv17@gmail.com" // Change this to the recipient's email address

	// Email content
	subject := "Server Shutdown Notification"
	body := "Your server has been shut down unexpectedly."

	// Constructing email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "base64"

	var msg bytes.Buffer
	for key, value := range headers {
		msg.WriteString(key + ": " + value + "\r\n")
	}
	msg.WriteString("\r\n" + base64.StdEncoding.EncodeToString([]byte(body)))

	// Connect to the SMTP server with TLS
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, msg.Bytes())
	if err != nil {
		log.Println("Failed to send shutdown notification email:", err)
		return
	}

	log.Println("Shutdown notification email sent successfully")
}
