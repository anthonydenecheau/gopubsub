package config

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mailjet/mailjet-apiv3-go"
)

const (
	format = time.RFC3339Nano
)

type MailAuthHook struct {
	AppName       string
	APIKeyPublic  string
	APIKeyPrivate string
	Subject       string
	Filename      *os.File
	From          *mail.Address
	To            *mail.Address
	levels        []logrus.Level
}

// NewMailAuthHook : send over MailJet
func NewMailAuthHook(appname string, apiKeyPublic string, apiKeyPrivate string, subject string, from string, to string, filename *os.File, levels []logrus.Level) (*MailAuthHook, error) {

	// Validate sender and recipient
	sender, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}
	receiver, err := mail.ParseAddress(to)
	if err != nil {
		return nil, err
	}

	return &MailAuthHook{
		AppName:       appname,
		APIKeyPublic:  apiKeyPublic,
		APIKeyPrivate: apiKeyPrivate,
		Subject:       subject,
		From:          sender,
		To:            receiver,
		Filename:      filename,
		levels:        levels}, nil
}

// Fire is called when a log event is fired.
func (hook *MailAuthHook) Fire(entry *logrus.Entry) error {

	// https://github.com/mailjet/mailjet-apiv3-go
	m := mailjet.NewMailjetClient(hook.APIKeyPublic, hook.APIKeyPrivate)

	//body := entry.Time.Format(format) + " - " + entry.Message
	subject := fmt.Sprintf(hook.Subject, hook.AppName)
	//fields, _ := json.MarshalIndent(entry.Data, "", "\t")
	//contents := fmt.Sprintf("Subject: %s\r\n\r\n%s\r\n\r\n%s", subject, body, fields)

	//fmt.Printf("filename: %s\n", hook.Filename)
	encodedLog := encodeBase64(hook.Filename)

	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: hook.From.Address,
				Name:  "Ne pas repondre",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: hook.To.Address,
				},
			},
			Subject:  subject,
			TextPart: "Hey ! Wake up ! Something's Wrong Happened !!!",
			HTMLPart: fmt.Sprintf("<b>Hey ! Wake up ! Something's Wrong Happened !!!</b><br>%s", createMessage(entry)),
			Attachments: &mailjet.AttachmentsV31{
				mailjet.AttachmentV31{
					ContentType:   "text/plain",
					Filename:      filepath.Base(hook.Filename.Name()),
					Base64Content: encodedLog,
				},
			},
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	res, err := m.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)

	return nil
}

// Levels returns the available logging levels.
func (hook *MailAuthHook) Levels() []logrus.Level {
	//return hook.levels
	return logrus.AllLevels[:hook.levels[0]+1]
}

func createMessage(entry *logrus.Entry) string {
	body := entry.Time.Format(format) + " - " + entry.Message
	fields, _ := json.MarshalIndent(entry.Data, "", "\t")
	contents := fmt.Sprintf("%s\r\n\r\n%s", body, fields)
	return contents
}

func encodeBase64(f *os.File) string {

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}
