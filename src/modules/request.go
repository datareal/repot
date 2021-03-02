package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Request is used in our lambda functions to define APIGateway Requests
type Request events.APIGatewayProxyRequest

// PostRequest is used to create the Request to the Slack Webhook url
func PostRequest(message string, SlackWebhook string) (string, int) {
	bytesRepresentation, err := json.Marshal(map[string]interface{}{
		"text": message,
	})
	if err != nil {
		fmt.Println("json:Marshal::Error ", err)
		return "Error when creating the Slack form", 500
	}

	resp, err := http.Post(SlackWebhook, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println("http:Post::Error: ", err)
		return "Error when making the POST request to the Slack Webhook URL", 500
	}

	return "Message created successfuly!", 200
}
