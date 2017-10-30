package main

import (
	"rider"
	"fmt"
	"rider/smtp/FlyWhisper"
	"errors"
)

func main() {
	app := rider.New()
	rlog := app.Logger(8)
	//wd, _ := os.Getwd()
	rlog.SetLogOutPath("")
	rlog.SmtpLogger(
		"postmaster@seemrice.com",
		"Hyp2Wxm2Hxy",
		"smtp.mxhichina.com",
		"25",
		"postmaster@seemrice.com",
	)
	//rlog.SetDestination(1)
	//rlog.AddDestination(0)
	fmt.Println(rlog.GetDestination())
	rlog.SetLogFileMaxSize(20 << 20)
	app.GET("/logger", func(c rider.Context) {
			c.SetHeader("ACCESS-CONTROL-ALLOW-ORIGIN", "*")
			rlog.INFO("xx", "yy")
			rlog.DEBUG("OK")
			//rlog.PANIC("adad")
			mess := FlyWhisper.NewMessage("logger send", []string{"1825909531@qq.com"})
			mess.AddHtml("<p>this is a mail logger</p>")
			c.Logger().SendMail(mess)
			c.SendJson(200, map[string]string{
				"a": "1",
			})
	})
	app.GET("/log500", func(c rider.Context) {
			c.Send(200, []byte("error500"))
	})
	app.GET("/300", func(c rider.Context) {
			c.Hijack()
			//panic(errors.New("adad"))
			c.Send(200, []byte("300"))
			c.Send(200, []byte("awd"))
			//c.SetStatusCode(400)
	})

	app.GET("/panic", func(context rider.Context) {
			panic(errors.New("adadad"))
	})
	app.Kid("/ada", &rider.Router{})
	app.Listen(":5000")

}
