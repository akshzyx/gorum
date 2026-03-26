package email

import (
	"fmt"
	"log"

	resend "github.com/resend/resend-go/v3"
)

type EmailSender struct {
	client  *resend.Client
	from    string
	baseURL string
}

func NewEmailSender(apiKey, from, baseURL string) *EmailSender {
	client := resend.NewClient(apiKey)

	return &EmailSender{
		client:  client,
		from:    from,
		baseURL: baseURL,
	}
}


func (s *EmailSender) SendVerificationEmail(to, token string) error {
	link := fmt.Sprintf("%s/auth/activate?token=%s", s.baseURL, token)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Activate your account",
		Html: fmt.Sprintf(`
			<h2>Verify your account</h2>
			<p>Click the link below to activate your account:</p>
			<a href="%s">Activate Account</a>
		`, link),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Println("[EMAIL ERROR]", err)
		return err
	}

	return nil
}
