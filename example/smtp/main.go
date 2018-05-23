package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hypwxm/rider/smtp/FlyWhisper"
)

func main() {
	mailer := FlyWhisper.NewMailer("postmaster@seemrice.com", "Hyp2Wxm2Hxy", "smtp.mxhichina.com", "25", "postmaster@seemrice.com")
	mess := FlyWhisper.NewMessage("test", []string{"1825909531@qq.com"})
	mess.AddText("<p>adawdawd</p>")
	mess.AddHtml("<p style='color:#ff0000;'>html<img src='http://img.51dbsd.com/Dealer/dc3130a9-4d27-4798-8fee-c6e5efd7e4c5.jpg'><p><a href='http://mv.jtwdxt.com'>qwdawdawdawdawd</a>")
	wd, _ := os.Getwd()
	mess.AddHtml("<div>123123</div>")
	mess.AddRcpt([]string{"1972650704@qq.com"})
	mess.RemoveRcpt("1825909531@qq.com")
	mess.SetHtml("<p style='font-size:30px;'>youjian</p>")
	//mess.AddInline(filepath.Join(wd, "src/github.com/hypwxm/rider/example/smtp/test.png"))
	//mess.AddInline(filepath.Join(wd, "src/github.com/hypwxm/rider/example/smtp/test.png"))
	mess.AddAttach(filepath.Join(wd, "src/github.com/hypwxm/rider/example/smtp/test.png"))
	//mess.AddAttach(filepath.Join(wd, "src/github.com/hypwxm/rider/example/smtp/54cc82980afd9.ai"))

	mess.RemoveHtml()
	mess.RemoveInline()
	mess.RemoveAttach()
	err := mailer.Send(mess)
	if err != nil {
		log.Println(err)
	}
}
