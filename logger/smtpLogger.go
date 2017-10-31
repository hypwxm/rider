package logger

import (
	"github.com/hypwxm/rider/smtp/FlyWhisper"
	"errors"
)

func (lq *LogQueue) SmtpLogger(username string, password string, host string, port string, from string) *FlyWhisper.SMTPSender {
	mail := FlyWhisper.NewMailer(username, password, host, port, from)
	lq.smtpLogger = mail
	lq.AddDestination(2)
	return mail
}

func (lq *LogQueue) SendMail(mess *FlyWhisper.Message) error {
	if !lq.DestExist(2) {
		return errors.New("mailLogger is now not available, if you need send log to email, please use method AddDestination(2)")
	}
	return lq.smtpLogger.Send(mess)
}