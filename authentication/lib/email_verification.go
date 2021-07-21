package authentication

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"strconv"
	"text/template"
)

func generateRandomCode() int64 {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return n.Int64()
}

func SendMail(email string) int {
	verificationCode := int(generateRandomCode())
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

	t, _ := template.ParseFiles("lib/email-verification-template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Email Verification \n%s\n\n", mimeHeaders)))

	err := t.Execute(&body, struct {
		VerificationCode string
	}{
		VerificationCode: strconv.Itoa(verificationCode),
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
