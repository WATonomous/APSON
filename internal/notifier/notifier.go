package notifier

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SendEmail sends an email notification with the given subject and body to all recipients.
func SendEmail(smtpServer string, smtpPort int, sender, password string, recipients []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	msg := buildMessage(sender, recipients, subject, body)
	auth := smtp.PlainAuth("", sender, password, smtpServer)
	return smtp.SendMail(addr, auth, sender, recipients, []byte(msg))
}

func buildMessage(sender string, recipients []string, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", sender),
		fmt.Sprintf("To: %s", strings.Join(recipients, ", ")),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
	}
	return strings.Join(headers, "\r\n") + body
}
