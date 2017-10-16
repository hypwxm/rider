package FlyWhisper

import "testing"

func TestMessage(t *testing.T) {
	mess := NewMessage("test", []string{"1@qq.com"})
	mess.AddRcpt([]string{"2@qq.com"})
	if len(mess.to) != 2 {
		t.Error("AddRcpt error")
	}
	if mess.to[0] != "1@qq.com" || mess.to[1] != "2@qq.com" {
		t.Error("AddRcpt error")
	}
	mess.SetRcpt([]string{"3@qq.com"})
	if len(mess.to) != 1 || mess.to[0] != "3@qq.com" {
		t.Error("SetRcpt error")
	}
	rc := mess.RemoveRcpt("3@qq.com")
	rc2 := mess.RemoveRcpt("xxx")
	if rc != "3@qq.com" || rc2 != "" || len(mess.to) != 0 {
		t.Error("RemoveRcpt error")
	}

	mess.AddInline("xxx.jpg")
	if len(mess.inlineImgs) != 1 || mess.inlineImgs[0].Filename != "xxx.jpg" || mess.inlineImgs[0].contentType != "image/jpeg" {
		t.Error("AddInlines error")
	}
	mess.SetInline("aaa.jpg")
	if len(mess.inlineImgs) != 1 || mess.inlineImgs[0].Filename != "aaa.jpg" || mess.inlineImgs[0].contentType != "image/jpeg" {
		t.Error("AddInlines error")
	}
	mess.AddInline("bbb.jpg")
	if len(mess.inlineImgs) != 2 || mess.inlineImgs[0].Filename != "aaa.jpg" || mess.inlineImgs[0].contentType != "image/jpeg" {
		t.Error("AddInlines error")
	}
	mess.RemoveInline()
	if len(mess.inlineImgs) != 0 {
		t.Error("RemoveInlines error")
	}


	mess.AddAttach("xxx.jpg")
	if len(mess.attachments) != 1 || mess.attachments[0].Filename != "xxx.jpg" || mess.attachments[0].contentType != "image/jpeg" {
		t.Error("AddAttachment error")
	}
	mess.SetAttach("aaa.jpg")
	if len(mess.attachments) != 1 || mess.attachments[0].Filename != "aaa.jpg" || mess.attachments[0].contentType != "image/jpeg" {
		t.Error("SetAttachment error")
	}
	mess.AddAttach("bbb.jpg")
	if len(mess.attachments) != 2 || mess.attachments[0].Filename != "aaa.jpg" || mess.attachments[0].contentType != "image/jpeg" {
		t.Error("AddAttachment error")
	}
	mess.RemoveAttach()
	if len(mess.attachments) != 0 {
		t.Error("RemoveAttachment error")
	}


	mess.AddText("aaa")
	if string(mess.text) != "aaa" {
		t.Error("AddText error")
	}
	mess.AddText("bbb")
	if string(mess.text) != "aaabbb" {
		t.Error("AddText error")
	}
	mess.SetText("ccc")
	if string(mess.text) != "ccc" {
		t.Error("SetText error")
	}
	mess.RemoveText()
	if string(mess.text) != "" {
		t.Error("RemoveText error")
	}


	mess.AddHtml("aaa")
	if string(mess.html) != "aaa" {
		t.Error("AddHtml error")
	}
	mess.AddHtml("bbb")
	if string(mess.html) != "aaabbb" {
		t.Error("AddHtml error")
	}
	mess.SetHtml("ccc")
	if string(mess.html) != "ccc" {
		t.Error("SetHtml error")
	}
	mess.RemoveHtml()
	if string(mess.html) != "" {
		t.Error("RemoveHtml error")
	}

}
