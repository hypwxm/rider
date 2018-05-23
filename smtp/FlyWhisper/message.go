package FlyWhisper

import (
	"errors"
	"mime"
	"path/filepath"
	"strings"

	"github.com/hypwxm/rider/utils/cryptos"
)

//定义发送的内容
type Message struct {
	subject     string
	text        []byte
	html        []byte
	inlineImgs  []*inlineImg
	attachments []*attachment
	to          []string
}

type inlineImg struct {
	Filename    string
	contentType string
	contentId   string
}

//上传附件
type attachment struct {
	Filename    string
	contentType string
	contentId   string
}

func NewMessage(subject string, to []string) *Message {
	return &Message{
		subject: subject,
		to:      to,
	}
}

//添加接收者
func (m *Message) AddRcpt(users []string) {
	m.to = append(m.to, users...)
}

//设置接收者
func (m *Message) SetRcpt(users []string) {
	m.to = users
}

//删除接收者
//删除成功返回删除的名称
//不存在返回""
func (m *Message) RemoveRcpt(user string) string {
	lenUsers := len(m.to)
	for k, name := range m.to {
		if name == user {
			if lenUsers-1 == k {
				m.to = m.to[:k]
			} else {
				m.to = append(m.to[:k], m.to[k+1:]...)
			}
			return name
		}
	}
	return ""
}

//添加inlineImg
func (m *Message) AddInline(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("Invalid filename (Empty string) ")
	}
	inline := &inlineImg{
		Filename: filename,
	}
	inline.contentType = mime.TypeByExtension(filepath.Ext(filename))
	inline.contentId = cryptos.RandString()
	m.inlineImgs = append(m.inlineImgs, inline)
	return nil
}

//设置inlineImg
func (m *Message) SetInline(filename string) error {
	m.inlineImgs = []*inlineImg{}
	return m.AddInline(filename)
}

//删除inlineImg
func (m *Message) RemoveInline() {
	m.inlineImgs = []*inlineImg{}
}

//添加attachment
func (m *Message) AddAttach(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("Invalid filename (Empty string) ")
	}
	attach := &attachment{
		Filename: filename,
	}
	attach.contentType = mime.TypeByExtension(filepath.Ext(filename))
	attach.contentId = cryptos.RandString()
	m.attachments = append(m.attachments, attach)
	return nil
}

//设置inlineImg
func (m *Message) SetAttach(filename string) error {
	m.attachments = []*attachment{}
	return m.AddAttach(filename)
}

//删除inlineImg
func (m *Message) RemoveAttach() {
	m.attachments = []*attachment{}
}

//添加text
func (m *Message) AddText(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	m.text = append(m.text, []byte(text)...)
}

//设置text
func (m *Message) SetText(text string) {
	m.text = []byte(text)
}

//删除text
func (m *Message) RemoveText() {
	m.text = []byte{}
}

//添加html
func (m *Message) AddHtml(html string) {
	if strings.TrimSpace(html) == "" {
		return
	}
	m.html = append(m.html, []byte(html)...)
}

//设置text
func (m *Message) SetHtml(html string) {
	m.html = []byte(html)
}

//删除text
func (m *Message) RemoveHtml() {
	m.html = []byte{}
}
