package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := r.Body
	fmt.Println(body)
	model := &model{}
	if err := json.Unmarshal([]byte(body), model); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	c, _ := json.Marshal(model)

	fmt.Println(string(c))
	task := model.Request.Intent.Slots.Utterance.Value
	user := model.Session.User.UserID
	dao := NewDAO(NewClient(http.DefaultClient))
	if err := dao.CreateTask(task, user); err != nil {
		task = err.Error()
	}
	outBody := `{
			"response": {
			  "outputSpeech": {
				"type": "PlainText",
				"text": "added task to ` + task + `"  
			  },
			  "shouldEndSession": true
			}
		  }`

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       outBody,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

type model struct {
	Request request `json:"request"`
	Session session `json:"session"`
}

type session struct {
	User user `json:"user"`
}

type user struct {
	UserID string `json:"userId"`
}

type request struct {
	Intent intent `json:"intent"`
}

type intent struct {
	Slots slots `json:"slots"`
}

type slots struct {
	Utterance utterance `json:"utterance"`
}

type utterance struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	ConfirmationStatus string `json:"confirmationStatus"`
	Source             string `json:"source"`
}
