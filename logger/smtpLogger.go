package logger

/*
	邮件日志
	SmtpLogger配置发送者信息
	SendMail在*FlyWhisper.Message中定义发送的内容，接受者等，发送出去
*/

import (
	"errors"
	"github.com/hypwxm/rider/smtp/FlyWhisper"
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
