package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SmileL1ne/web-mailing/db"
	"github.com/SmileL1ne/web-mailing/model"
	"github.com/SmileL1ne/web-mailing/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailApp struct {
	MailDB   db.DataStore
	MailChan chan model.Mail
}

func NewMailApp(client *mongo.Client, mailChan chan model.Mail) Logic {
	return &MailApp{
		MailDB:   db.NewMongo(client),
		MailChan: mailChan,
	}
}

func (ma *MailApp) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tools.HTMLRender(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (ma *MailApp) GetSubscriber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sub model.Subscriber
		subsriber, err := tools.ReadForm(r, sub)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read json: %v\n", err), http.StatusBadRequest)
			return
		}
		ok, msg, err := ma.MailDB.AddSubscriber(subsriber)
		if err != nil {
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		switch ok {
		case msg == "":
			tools.JSONWriter(w, "You have already registered", http.StatusOK)
		case msg != "":
			tools.JSONWriter(w, msg, http.StatusOK)
		}
	}
}

func (ma *MailApp) SendMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mailUpload model.MailUpload
		upload, err := tools.ReadMultiForm(w, r, mailUpload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg, err := ma.MailDB.AddMail(upload)
		if err != nil {
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		log.Println(msg)
		log.Println("........ Preparing send mail to subscribers ........")
		time.Sleep(time.Millisecond)
		log.Println("........ Accessing the Subscribers database ........")

		res, err := ma.MailDB.FindSubscribers()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed query: %v\n", err), http.StatusInternalServerError)
			return
		}
		for _, s := range res {
			subEmail := s["email"].(string)
			firstName := s["first_name"].(string)
			lastName := s["last_name"].(string)
			subName := fmt.Sprintf("%s %s", firstName, lastName)

			mail := model.Mail{
				Source:      os.Getenv("GMAIL_ACC"),
				Destination: subEmail,
				Message:     upload.DocxContent,
				Subject:     upload.DocxName,
				Name:        subName,
			}
			ma.MailChan <- mail
		}

		err = tools.JSONWriter(w, fmt.Sprintf("Mail sent to %d subscribers", len(res)), http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (ma *MailApp) DeleteSubscriber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read from unsubscribe form: %v", err), http.StatusBadGateway)
			return
		}
		email := r.Form.Get("email")
		ok, err := ma.MailDB.DeleteSubscriber(email)
		if !ok {
			err = tools.JSONWriter(w, "No such subscriber", http.StatusBadRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete subscriber: %v", err), http.StatusInternalServerError)
		}
		err = tools.JSONWriter(w, "Sucessfully deleted subscriber", http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
