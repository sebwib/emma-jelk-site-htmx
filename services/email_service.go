package services

import (
	"fmt"
	"os"

	"github.com/sebwib/emma-site-htmx/db"
	gomail "gopkg.in/mail.v2"
)

func sendEmail(subject string, body string) error {
	googleAppPassword := os.Getenv("GOOGLE_APP_PASSWORD")
	fromAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	recipientAddress := os.Getenv("EMAIL_RECIPIENT_ADDRESS")

	if fromAddress == "" || recipientAddress == "" || googleAppPassword == "" {
		return fmt.Errorf("email configuration is missing")
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", fromAddress)
	message.SetHeader("To", recipientAddress)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, fromAddress, googleAppPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}

func SendOrder(buyerEmail string, order db.Order) error {
	subject := "New Order Received"
	body := "You have received a new order:\n\n"
	body += "Buyer Email: " + buyerEmail + "\n"
	body += "Order Details:\n"
	for _, item := range order.Rows {
		body += fmt.Sprintf("- Print ID: %s, Type: %s, Quantity: %d, Price per unit: %.2f\n", item.Title, item.Typ, item.Quantity, item.Price)
	}

	return sendEmail(subject, body)
}
