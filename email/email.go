package email

import (
	"log"
	"os"

	"github.com/SmileL1ne/web-mailing/model"
	"gopkg.in/gomail.v2"
)

func MailServer(mailChan model.Mail) {
	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("USER_NAME"), os.Getenv("APP_PASSWORD"))
	s, err := d.Dial()
	if err != nil {
		log.Panicf("Error connecting to the Mail Server: %v\n", err)
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", mailChan.Source, os.Getenv("USER_NAME"))
	msg.SetHeader("To", mailChan.Destination)
	msg.SetHeader("Subject", mailChan.Subject)
	msg.SetBody("text/html", mailChan.Message)

	if err := gomail.Send(s, msg); err != nil {
		log.Printf("Mail Server: %s %v\n", mailChan.Destination, err)
	}
	msg.Reset()
}

func MailDelivery(mailChan <-chan model.Mail, worker int) {
	completionChan := make(chan bool, worker)
	for i := 0; i < worker; i++ {
		go func() {
			defer func() {
				completionChan <- true
			}()
			for m := range mailChan {
				MailServer(m)
			}
		}()
	}
	for i := 0; i < worker; i++ {
		<-completionChan
	}
}
