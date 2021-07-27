package authentication

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"text/template"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func generateRandomCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SendMail(email string) string {
	verificationCode := generateRandomCode(6)
	// Sender data.
	from := "mineloopeducation@gmail.com"
	password := "Hungthjkju2"

	//Recipient
	////////////////////// Remember to Change
	to := []string{"blankspace171098@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	//message := []byte("This is a test email message.")

	//Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	///////// REMEMBER RUN AS SERVER DIRECTORY
	t, _ := template.ParseFiles("lib/email-verification-template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Email Verification \n%s\n\n", mimeHeaders)))

	err := t.Execute(&body, struct {
		VerificationCode string
	}{
		VerificationCode: verificationCode,
	})
	if err != nil {
		log.Fatalln(err)
	}
	//Sending email.
	err2 := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err2 != nil {
		log.Fatalln(err2)
	}
	return verificationCode
}
