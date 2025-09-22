package service

import "speakpall/pkg/mailer"

type MailerService interface {
	Send(to, subject, body string) error
}

type mailerService struct {
	mailer *mailer.Mailer
}

func NewMailerService(m *mailer.Mailer) MailerService {
	return &mailerService{mailer: m}
}

func (m *mailerService) Send(to, subject, body string) error {
	return m.mailer.Send(to, subject, body)
}
