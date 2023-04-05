package smtpService

import (
	smtpDomain "go-rotate-backup-s3/commons/domain/smtp"
	gomailInfra "go-rotate-backup-s3/commons/infra/gomail"
)

type SmtpService struct {
	mailer gomailInfra.ISmtp
}

type ISmtpService interface {
	Send(smtpDomain.SendArgs) error
}

func New(mailer gomailInfra.ISmtp) ISmtpService {
	return &SmtpService{mailer: mailer}

}

func (s *SmtpService) Send(args smtpDomain.SendArgs) error {
	err := s.mailer.Send(args)
	return err
}
