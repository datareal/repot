package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/datareal/repot/src/modules"
)

type message struct {
	Message string `json:"message"`
}

// Handler is used by lambda start which start our code
func Handler(request modules.Request) (modules.Response, error) {
	// Verify if the request.Body is base64
	// if it is return an error requesting to add
	// 'content-type': 'application/json' to the header
	_, decodeErr := base64.StdEncoding.DecodeString(request.Body)
	if decodeErr == nil {
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Failed to create the Request body, try setting 'Content-Type': 'application/json' on your request header.",
		}, 500)
	}

	var Message message
	err := json.Unmarshal([]byte(request.Body), &Message)
	if err != nil {
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Error when reading request data",
		}, http.StatusInternalServerError)
	}

	SlackWebhook, ok := os.LookupEnv("SLACK_WEBHOOK_REPORTS")
	if !ok {
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Failed loading Slack Webhook URL",
		}, 500)
	}
	fmt.Println("SlackWebhook URL: ", SlackWebhook)

	response, statusCode := modules.PostRequest(Message.Message, SlackWebhook)

	return modules.CreateResponse(modules.ResponseMessage{
		Message: response,
	}, statusCode)
}

func main() {
	lambda.Start(Handler)
}
