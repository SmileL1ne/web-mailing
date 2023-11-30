package model

import "time"

// Information from subscribers
type Subscriber struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
	Interest  string `bson:"interest" json:"interest"`
}

// Holds content and other info for the mail
type MailUpload struct {
	DocxName    string    `bson:"docx_name" json:"docx_name"`
	DocxContent string    `bson:"docx" json:"docx"`
	Date        time.Time `bson:"date" json:"date"`
}

// Holds mail's properties
type Mail struct {
	Source      string
	Destination string
	Message     string
	Subject     string
	Name        string
}
