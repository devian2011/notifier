package transports

import (
	"fmt"
	"net/smtp"

	"notifications/internal/dto"
)

type SmtpConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	User     string `yaml:"user" yaml:"user"`
	Password string `yaml:"password" yaml:"password"`
	From     string `yaml:"from" yaml:"from"`
}

type SmtpSender struct {
	cfg *SmtpConfig
}

func (s *SmtpSender) Send(to []string, msg *dto.Message, meta map[string]string) error {
	smtpMessage := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", s.cfg.From, msg.Subject, msg.Body)

	return smtp.SendMail(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host),
		s.cfg.User, to, []byte(smtpMessage))
}
