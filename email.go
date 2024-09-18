package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"path/filepath"
)

func sendEmailWithGmail(to, subject, body, attachment string) error {
	from := os.Getenv("GMAIL_USERNAME")
	password := os.Getenv("GMAIL_PASSWORD")

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Create the email body and headers
	var message bytes.Buffer
	boundary := "my-boundary-12345"

	// Write headers
	message.WriteString(fmt.Sprintf("From: %s\n", from))
	message.WriteString(fmt.Sprintf("To: %s\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	message.WriteString("MIME-Version: 1.0\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	message.WriteString("\n")

	// Write body part
	message.WriteString(fmt.Sprintf("--%s\n", boundary))
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\n")
	message.WriteString("\n")
	qp := quotedprintable.NewWriter(&message)
	qp.Write([]byte(body))
	qp.Close()
	message.WriteString("\n")

	// Add attachment if provided
	if attachment != "" {
		fileData, err := ioutil.ReadFile(attachment)
		if err != nil {
			return fmt.Errorf("failed to read attachment: %v", err)
		}

		fileName := filepath.Base(attachment)
		encodedFile := base64.StdEncoding.EncodeToString(fileData)

		// Write attachment part
		message.WriteString(fmt.Sprintf("--%s\n", boundary))
		message.WriteString("Content-Type: application/octet-stream\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n", fileName))
		message.WriteString("Content-Transfer-Encoding: base64\n")
		message.WriteString("\n")
		message.WriteString(encodedFile)
		message.WriteString("\n")
	}

	// End the MIME message
	message.WriteString(fmt.Sprintf("--%s--\n", boundary))

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Println("Email sent successfully with attachment")
	return nil
}
