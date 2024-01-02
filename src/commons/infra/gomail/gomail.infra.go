package gomailInfra

import (
	"bytes"
	"crypto/tls"
	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	smtpDomain "go-rotate-backup-s3/commons/domain/smtp"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/gomail.v2"
)

type ISmtp interface {
	Send(smtpDomain.SendArgs) error
}

type GomailInfra struct {
	Smtp *gomail.Dialer
}

func New() ISmtp {
	d := gomail.NewDialer(
		configService.GetSmtpHost(),
		configService.GetSmtpPort(),
		configService.GetSmtpUser(),
		configService.GetSmtpPass(),
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &GomailInfra{Smtp: d}
}

func (s *GomailInfra) Send(args smtpDomain.SendArgs) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", configService.GetSmtpFrom())
	msg.SetHeader("To", args.To)
	if args.Cc != "" {
		msg.SetAddressHeader("Cc", args.Cc, "")
	}

	msg.SetHeader("Subject", args.Subject)

	pathTemplate, err := filepath.Abs(configService.GetSmtpTemplate() + args.Template)

	if err != nil {
		logService.Error(err.Error())
		os.Exit(1)
	}

	templateEmail, errParse := template.ParseFiles(pathTemplate)

	if errParse != nil {
		logService.Error(errParse.Error())
		return errParse
	}

	var html bytes.Buffer

	templateEmail.Execute(&html, args.Data)

	msg.SetBody("text/html", string(html.Bytes()))

	errSmtp := s.Smtp.DialAndSend(msg)

	if errSmtp != nil {
		logService.Error(errSmtp.Error())
	}

	logService.Info("email sent successful: " + args.To)

	return errSmtp

}
