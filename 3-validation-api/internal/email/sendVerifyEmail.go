package email

import (
	"net/smtp"
	"strings"
	"validationApi/configs"
	"validationApi/pkg/hash"
)

func SendVerifyEmail(emailTo []string, conf *configs.Config) error {
	from := conf.Email
	password := conf.Password

	to := emailTo

	smtpHost := conf.Address
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	hashStr := hash.Generate()
	message := []byte("verify http://localhost:8081/verify/" + hashStr)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	emailData := &EmailVerifyData{
		Email: strings.Join(to, ""),
		Hash:  hashStr,
	}

	err = SaveVerifyEmailData(*emailData)

	if err != nil {
		return err
	}

	return nil
}
