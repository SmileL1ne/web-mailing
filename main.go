package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SmileL1ne/web-mailing/db"
	"github.com/SmileL1ne/web-mailing/email"
	"github.com/SmileL1ne/web-mailing/handlers"
	"github.com/SmileL1ne/web-mailing/model"
	"github.com/joho/godotenv"
)

var (
	MailChan   chan model.Mail
	BufferSize int
	Worker     int
)

func main() {
	MailChan = make(chan model.Mail, BufferSize)
	Worker = 5

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting the Mail Server")

	log.Println("Preparing Database Connection")

	client := db.OpenConnect()
	defer func(ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(context.TODO())

	go email.MailDelivery(MailChan, Worker)
	defer close(MailChan)

	app := handlers.NewMailApp(client, MailChan)
	handle := Routes(app)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: handle,
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln("Shutting down the Mail App Server")
	}
}
