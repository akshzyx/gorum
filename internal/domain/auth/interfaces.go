package auth

type EmailSender interface {
	SendVerificationEmail(to, token string) error
}
