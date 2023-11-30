package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetConnect(uri string) (*mongo.Client, error) {
	dbCtx, dbCancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer dbCancelCtx()

	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicln("Error connecting to the database:", err)
	}

	if err := client.Ping(dbCtx, nil); err != nil {
		log.Fatalln("Cannot ping the database:", err)
	}

	return client, nil
}

func OpenConnect() *mongo.Client {
	uri := os.Getenv("URI")
	count := 0
	log.Println("......... Setting Up Mongo DB Connection .........")
	for {
		client, err := SetConnect(uri)
		if err != nil {
			log.Println("Connection to app-mail is not established")
			count++
		} else {
			log.Println("Connection to app-mail is established")
			return client
		}

		if count >= 5 {
			log.Println(err)
			return nil
		}

		log.Println("Wait: ... Retrying to connect to the mail-app ...")
		time.Sleep(10 * time.Second)
		continue
	}
}
