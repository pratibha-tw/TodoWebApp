package emailservice

import (
	"fmt"
	"net/smtp"
	"todoapp/config"
)

type EmailService interface {
	SendNotification(toEmail string, body []byte)
}

type emailService struct {
	config.EmailConfig
	auth smtp.Auth
}

func (e emailService) SendNotification(toEmail string, body []byte) {
	err := smtp.SendMail(e.Host+":"+e.Port, e.auth, e.FromEmail, []string{toEmail}, body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Email sent successfully!")
}

func (e *emailService) GetAuth(email_cfg config.EmailConfig) {
	auth := smtp.PlainAuth("", email_cfg.FromEmail, email_cfg.Password, email_cfg.Host)
	e.auth = auth
}
func NewEmailService(cfg config.EmailConfig) EmailService {
	var emailSer emailService
	emailSer.GetAuth(cfg)
	emailSer.EmailConfig = cfg
	fmt.Println("got email config")
	return &emailSer
}
