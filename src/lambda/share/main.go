package main

import (
	"os"
	"fmt"
	"bytes"
	"encoding/json"
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type BodyRequest struct {
	RequestName string `json:"name"`
}

func CreateResponse(message string, responseStatusCode int) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"Message": message,
	})
	if err != nil {
		return Response{StatusCode: 500, Body: "Failed to create the Request body"}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      responseStatusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}

func Handler(request Request) (Response, error) {
	// Verify if the request.Body is base64 
	// if it is return an error requesting to add
	// 'content-type': 'application/json' to the header
	_, decodeErr := base64.StdEncoding.DecodeString(request.Body)
	if decodeErr == nil {
		return CreateResponse("Failed to create the Request body, try setting 'Content-Type': 'application/json' on your request header.", 500)
	}

	_, ok := os.LookupEnv("SLACK_WEEBHOOK_REPORTS")
	if !ok {
		return CreateResponse("Failed loading Slack Webhook URL", 500)
	}

	fmt.Println("Received body: ", request.Body)

	return CreateResponse("Message sent successfuly!", 200)
}

func main() {
	lambda.Start(Handler)
}