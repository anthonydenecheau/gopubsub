package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	logger "github.com/anthonydenecheau/gopubsub/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type configLogger struct {
	Log  *logrus.Logger
	File *os.File
}

// NewLogger contructor
func NewLogger(config logger.LoggerConfiguration, app logger.AppConfiguration) (*logrus.Logger, *os.File) {
	log := logrus.New()
	logrus.SetLevel(logrus.InfoLevel)

	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = time.RFC3339Nano
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	if config.Hooks.Name == "file" {
		filename := fmt.Sprintf("%s/%s", config.Hooks.Options.Directory, config.Hooks.Options.Filename)
		logfile := fmt.Sprintf("%s.%s", filename, time.Now().Format("200601020000"))

		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		logf, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", filename, "%Y%m%d%H%M"),
			// KO ON WINDOWS ! rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(time.Duration(config.Hooks.Options.Maxday*24)*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			log.Printf("failed to create rotatelogs: %s", err)
			//log.Fatal(err)
		}
		log.SetOutput(logf)

		// Gestion du niveau d'alerte minimum pour l'envoi de mail
		var levels []logrus.Level
		levels = append(levels, logrus.ErrorLevel)
		//levels = append(levels, logrus.PanicLevel)
		//levels = append(levels, logrus.FatalLevel)

		// Logrus has seven logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.
		log.SetLevel(logrus.InfoLevel)

		// Authentication MailJet
		config.Hooks.Mail.MjAPIPublic = mustGetenv("MJ_APIKEY_PUBLIC")
		config.Hooks.Mail.MjAPIPrivate = mustGetenv("MJ_APIKEY_PRIVATE")

		hook, err := NewMailAuthHook(app.Name, config.Hooks.Mail.MjAPIPublic, config.Hooks.Mail.MjAPIPrivate, config.Hooks.Mail.Subject, config.Hooks.Mail.Sender, config.Hooks.Mail.Receivers, file, levels)
		if err == nil {
			log.Hooks.Add(hook)
		}

		return log, file
	}

	return log, nil
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
