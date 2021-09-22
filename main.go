package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jpastorm/dynamodb/repository"
)

func main() {
	lambda.Start(save)
}

var Posts []repository.Post

func save(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,HEAD,OPTIONS,POST",
		},
	}

	var post repository.Post
	err := json.Unmarshal([]byte(ev.Body), &post)
	if err != nil {
		panic(err)
	}

	repo := repository.NewDynamoDBRepository()
	_, err = repo.Save(&repository.Post{ID: time.Now().Unix(), Title: post.Title, Text: post.Text})
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp.StatusCode = 200
	resp.Body = fmt.Sprintln("Guardado con exito: ", post.Title)
	return resp, nil
}

func list(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := repository.NewDynamoDBRepository()
	posts, err := repo.FindAll()
	if err != nil {
		log.Println(err)
	}
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,HEAD,OPTIONS,POST",
		},
	}

	out, err := json.Marshal(posts)
	if err != nil {
		panic(err)
	}

	resp.Body = fmt.Sprintln(string(out))
	return resp, nil
}
