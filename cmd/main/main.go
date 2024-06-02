package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"notifications/internal"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT)
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Errorf("critical application err")
			stop()
		}
	}()

	logrus.Infoln("application init")
	app, initAppErr := internal.NewApp(ctx)
	if initAppErr != nil {
		logrus.WithField("error", initAppErr).Errorf("error on application init")
		return
	}
	logrus.Infoln("application init complete")
	logrus.Infoln("application running...")
	if execErr := app.Run(); execErr != nil {
		logrus.WithField("error", initAppErr).Errorf("error on application run")
		return
	}
}
