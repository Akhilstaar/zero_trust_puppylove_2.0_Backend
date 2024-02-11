package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

// TODO: Need to fix this code to authenticate only once & use the session for sending multiple emails.
func sendEmail(from, password, smtpHost, smtpPort, to string, bodyContent []byte) error {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, bodyContent)
	return err
}

func SendMail(name string, to string, authCode string) error {

	// smtp server configuration.
	from := os.Getenv("EMAIL_ID")
	password := os.Getenv("EMAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("From:PClub IITK <pclubiitk@gmail.com>\nSubject: Zero-trust Puppy Love Authentication Code \n%s\n\n", mimeHeaders)))

	mailTemplate := fmt.Sprintf("<html><style>.container{font-size: large;}body{font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;}code{padding: 6px;background-color: antiquewhite;border-radius: 5px;}</style><body><div class='container'>Hello %s,\n\nThis is your authentication code for Zero-trust-PuppyLove.\n\n <code>%s</code> \n\n .</div></body></html>", name, authCode)
	body.Write([]byte(mailTemplate))

	// Sending email
	err := sendEmail(from, password, smtpHost, smtpPort, to, body.Bytes())

	if err != nil {
		fmt.Printf("Error sending email to %s using the first email configuration: %v\n", to, err)

		// Second email configuration.
		from = os.Getenv("EMAIL_ID_2")
		password = os.Getenv("EMAIL_PASS_2")
		smtpHost = os.Getenv("SMTP_HOST_2")
		smtpPort = os.Getenv("SMTP_PORT_2")

		// Retry with the second email configuration.
		err = sendEmail(from, password, smtpHost, smtpPort, to, body.Bytes())
		if err != nil {
			fmt.Printf("Error sending email to %s using the second email configuration: %v\n", to, err)
			return err
		}
	}
	return nil
}
