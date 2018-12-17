package main

import (
	"flag"
	"fmt"
	"time"

	"tweg"

	"github.com/sirupsen/logrus"
)

var (
	tweet   = flag.String("tweet", "A koala arrives in the great forest of Wumpalumpa", "Tweet ?")
	secret  = flag.String("secret", "alpaga", "Secret ?")
	action  = flag.String("action", "encode", "Action ?")
	logging = flag.String("logging", "info", "Logging level")
)

func init() {
	// Parse the flags
	flag.Parse()

	// Set localtime to UTC
	time.Local = time.UTC

	// Set the logging level
	level, err := logrus.ParseLevel(*logging)
	if err != nil {
		logrus.Fatalln("Invalid log level ! (panic, fatal, error, warn, info, debug)")
	}
	logrus.SetLevel(level)

	// Set the TextFormatter
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
}

func main() {
	t := tweg.NewTweg()

	switch *action {
	case "encode":
		result, err := t.Encode(*tweet, *secret)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		logrus.Infoln("Result is :", result)
		break
	case "decode":
		fmt.Println("Result is :", t.Decode("A kｏａla arrivｅs іn the great forest of Wumpalumpa"))
		break
	default:
		logrus.Fatalln("Action is incorrect")
	}

}
