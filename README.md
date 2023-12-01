# Mailing Web App

## Overview

It is web application where you can subscribe to the mail and send word or txt document to all the subscribers. Acceptable formats - .doc, .docx, .txt

## Functionality

- Add subscriber
- Send document to all subscribers
- Delete subscriber

## Installation

> To start, at first **fork** this repository

### Environment setup

> In this project .env is used to keep all sensitive data out safe. There is .env.example file that shows variables needed to create in your .env file in order to run this application

1. Create .env file in root directory of project by given prototype (".env.example")
2. Fill this file with all the necessery information

> [!WARNING]
> Remember that .env file contains sensitive data that should shared openly

### Run application

1. Go to the root directory of the project
2. Run this command:
```
go run main.go routes.go
```
3. Local server has started on the port that you've written in your .env file - `http://localhost:<PORT>` 
