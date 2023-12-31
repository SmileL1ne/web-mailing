package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SmileL1ne/web-mailing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	MailDB *mongo.Client
}

func NewMongo(client *mongo.Client) DataStore {
	return &Mongo{MailDB: client}
}

func (mg *Mongo) AddSubscriber(subs model.Subscriber) (bool, string, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	var res bson.M
	filter := bson.D{{Key: "email", Value: subs.Email}}
	err := Default(mg.MailDB, "subscribers").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := Default(mg.MailDB, "subscribers").InsertOne(ctx, subs)
			if err != nil {
				return false, "", fmt.Errorf("AddSubscriber: cannot register this account: %v", err)
			}
			return true, fmt.Sprintf("New subscriber added"), nil
		}
		log.Fatalf("AddSubsriber: cannot query database: %v", err)
	}
	return true, "", nil
}

func (mg *Mongo) AddMail(mu model.MailUpload) (string, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	_, err := Default(mg.MailDB, "mails").InsertOne(ctx, mu)
	if err != nil {
		return "", fmt.Errorf("AddMail: unable to insert mail: %v\n", err)
	}
	return "New mail successfully uploaded", nil
}

func (mg *Mongo) FindSubscribers() ([]primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	var res []bson.M
	cursor, err := Default(mg.MailDB, "subscribers").Find(ctx, bson.D{})
	if err != nil {
		return []bson.M{}, err
	}
	if err := cursor.All(ctx, &res); err != nil {
		return nil, fmt.Errorf("FindSubscribers: cannot get all mails: %v\n", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("FindSubsribers: cursor error: %v\n", err)
	}
	return res, nil
}

func (mg *Mongo) DeleteSubscriber(email string) (bool, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	var res bson.M
	filter := bson.D{{Key: "email", Value: email}}
	err := Default(mg.MailDB, "subscribers").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, fmt.Errorf("Couldn't find subscriber - %s", email)
		}
		log.Fatalf("DeleteSubscriber: couldn't query database - %v", err)
	}
	_, err = Default(mg.MailDB, "subscribers").DeleteOne(ctx, filter)
	if err != nil {
		return true, fmt.Errorf("DeleteSubscriber: cannot delete subscriber - %v", err)
	}
	return true, nil
}
