package modules

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

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