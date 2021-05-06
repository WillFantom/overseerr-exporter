package main

import (
	"github.com/sirupsen/logrus"
	"github.com/willfantom/overseerr-exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if logrus.GetLevel() < logrus.WarnLevel {

		}
		logrus.WithField("extended", err.Error()).
			Fatalln("an error occurred executing the command")
	}
}
