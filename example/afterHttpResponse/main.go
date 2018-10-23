package main

import (
	"os"
	"time"

	"github.com/hypwxm/rider"
	"github.com/sirupsen/logrus"
)

func main() {
	app := rider.New()
	log := logrus.New()
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})
	app.AfterHttpResponse(func(ctx rider.Context, statusCode int, timeTaken time.Duration) {
		log.WithFields(logrus.Fields{
			"statusCode": statusCode,
			"timeTaken":  timeTaken.String(),
			"method":     ctx.Method(),
			"uid":        ctx.RequestID(),
			"uri":        ctx.Request().RequestURI(),
		}).Info()
	})
	app.GET("/", func(c rider.Context) {
		c.SendString(200, "xxxx")
	})
	app.Listen(":5000")
}
