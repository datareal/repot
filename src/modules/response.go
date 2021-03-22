package modules

import (
	"bytes"
	"encoding/json"
)

// Response is used in our lambda functions to define APIGateway Response
type Response struct {
	StatusCode int               `json:"StatusCode"`
	Body       string            `json:"Body"`
	Headers    map[string]string `json:"Headers"`
}

// ResponseMessage is used in our lambda function to return when there's a message to be returned
type ResponseMessage struct {
	Message string
}

// CreateResponse is used to create a response to the received request
func CreateResponse(message ResponseMessage, responseStatusCode int) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(message)
	if err != nil {
		return Response{StatusCode: 500, Body: "Failed to create the Request body"}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: responseStatusCode,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}
