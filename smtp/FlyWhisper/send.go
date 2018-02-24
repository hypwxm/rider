package FlyWhisper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/hypwxm/rider/utils/cryptos"
	"io/ioutil"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
)

//img suffix
var (
	ImgSuffix = []string{".bmp", ".jpg", ".jpeg", ".png", ".tiff", ".gif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".wmf"}
)

//mixed
//	alternative
//		text
//		related
//			html
//			inline image
//			inline image
//	attachment
//	attachment

//发送者
type SMTPSender struct {
	Username string
	Password string
	Host     string
	Port     string
	Header   textproto.MIMEHeader
	Subject  string
	From     string
	Text     string
	Html     string
}

func NewMailer(username string, password string, host string, port string, from string) *SMTPSender {
	return &SMTPSender{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		From:     from,
	}
}

func (s *SMTPSender) getAuth() smtp.Auth {
	if len(strings.Trim(s.Username, " ")) == 0 && len(strings.Trim(s.Password, " ")) == 0 {
		return nil
	}
	return smtp.PlainAuth(
		"",
		s.Username,
		s.Password,
		s.Host,
	)
}

func (s *SMTPSender) mailByte(message *Message) ([]byte, error) {
	var buff *bytes.Buffer = &bytes.Buffer{}
	if len(message.to) == 0 {
		return nil, errors.New("you have not been set receiver")
	}
	buff.WriteString("To: " + strings.Join(message.to, ";") + "\r\n")
	if len(s.From) == 0 {
		return nil, errors.New("you have not been set sender")
	}
	var (
		boundaryMixed       string = "@#riderMixed#@"
		boundaryRelated     string = "@#riderRelated#@"
		boundaryAlternative string = "@#riderAlternative#@"
	)
	buff.WriteString("From:" + s.Username + " <" + s.From + ">\r\n")
	if strings.TrimSpace(message.subject) != "" {
		buff.WriteString("Subject: " + message.subject + "\r\n")
	} else {
		buff.WriteString("Subject: " + s.Subject + "\r\n")
	}
	buff.WriteString("MESSAGE-Id: " + cryptos.RandString() + "\r\n")
	buff.WriteString("MIME-Version: 1.0\r\n")
	buff.WriteString("Content-Type: multipart/mixed; boundary=" + boundaryMixed + "\r\n\r\n")
	buff.WriteString("--" + boundaryMixed + "\r\n")
	buff.WriteString("Content-Type: multipart/related; boundary=" + boundaryRelated + "\r\n\r\n")
	buff.WriteString("--" + boundaryRelated + "\r\n")
	buff.WriteString("Content-Type: multipart/alternative; boundary=" + boundaryAlternative + "\r\n\r\n")
	buff.WriteString("--" + boundaryAlternative + "\r\n")
	if message == nil {
		return nil, errors.New("nil message is not allowed")
	}
	//start text/plain
	if strings.TrimSpace(s.Text) != "" || strings.TrimSpace(string(message.text)) != "" {
		buff.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buff.WriteString("Content-Transfer-Encoding: " + strconv.Itoa(len(s.Text)+len(message.text)) + "bit\r\n\r\n")
		if strings.TrimSpace(s.Text) != "" {
			//global text ,this will contain every mail if was defined
			buff.WriteString(s.Text)
		}
		if strings.TrimSpace(string(message.text)) != "" {
			//submessage text
			buff.Write(message.text)
		}
	}
	buff.WriteString("\r\n\r\n--" + boundaryAlternative + "\r\n")

	//start text/html
	if strings.TrimSpace(s.Html) != "" || strings.TrimSpace(string(message.html)) != "" {
		buff.WriteString("Content-Type: text/html; charset=utf-8\r\n")
		buff.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		if strings.TrimSpace(s.Html) != "" {
			//global text ,this will contain every mail if was defined
			buff.WriteString(s.Html)
		}
		if strings.TrimSpace(string(message.html)) != "" {
			//submessage text
			buff.Write(message.html)
		}
	}
	buff.WriteString("\r\n\r\n--" + boundaryAlternative + "\r\n")

	//start inlineImgs
	if len(message.inlineImgs) != 0 {
		for _, img := range message.inlineImgs {
			if !isImage(img.Filename) {
				continue
			}
			imgData, err := ioutil.ReadFile(img.Filename)
			if err != nil {
				return nil, err
			}
			buff.WriteString("\r\n--" + boundaryRelated + "\r\n")
			buff.WriteString("Content-Type: " + img.contentType + "\r\n")
			buff.WriteString("Content-Transfer-Encoding: base64\r\n")
			buff.WriteString("Content-Disposition: inline; filename=\"" + filepath.Base(img.Filename) + "\"\r\n\r\n")
			base64Data := make([]byte, base64.StdEncoding.EncodedLen(len(imgData)))
			base64.StdEncoding.Encode(base64Data, imgData)
			buff.Write(base64Data)
			buff.WriteString("\r\n\r\n--" + boundaryRelated + "\r\n")
		}
	}

	//start attachment
	if len(message.attachments) != 0 {
		for _, attach := range message.attachments {
			attachData, err := ioutil.ReadFile(attach.Filename)
			if err != nil {
				return nil, err
			}
			buff.WriteString("\r\n--" + boundaryMixed + "\r\n")
			buff.WriteString("Content-Type: " + attach.contentType + "\r\n")
			buff.WriteString("Content-Transfer-Encoding: base64\r\n")
			buff.WriteString("Content-Disposition: attachment; filename=\"" + filepath.Base(attach.Filename) + "\"\r\n\r\n")

			base64Data := make([]byte, base64.StdEncoding.EncodedLen(len(attachData)))
			base64.StdEncoding.Encode(base64Data, attachData)
			buff.Write(base64Data)
			buff.WriteString("\r\n\r\n--" + boundaryMixed + "\r\n")
		}
	}
	return buff.Bytes(), nil
}

func (s *SMTPSender) Send(message *Message) error {
	if message == nil {
		return errors.New("message can not be nil")
	}
	buffByte, err := s.mailByte(message)
	//fmt.Printf("%s", buffByte)
	if err != nil {
		return err
	}
	err = smtp.SendMail(s.Host+":"+s.Port, s.getAuth(), s.From, message.to, buffByte)
	if err != nil {
		return err
	}
	return nil
}

//判断是否为图片
func isImage(filename string) bool {
	suffix := filepath.Ext(filename)
	allextstr := strings.Join(ImgSuffix, ",")
	if strings.Index(allextstr, suffix) > -1 {
		return true
	}
	return false
}
