package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/datareal/repot/src/modules"
)

// Handler is used by lambda start which start our code
func Handler(request modules.Request) (modules.Response, error) {
	// Verify if the request.Body is base64
	// if it is return an error requesting to add
	// 'content-type': 'application/json' to the header
	_, decodeErr := base64.StdEncoding.DecodeString(request.Body)
	if decodeErr == nil {
		return modules.CreateResponse("Failed to create the Request body, try setting 'Content-Type': 'application/json' on your request header.", 500)
	}

	_, ok := os.LookupEnv("SLACK_WEEBHOOK_REPORTS")
	if !ok {
		return modules.CreateResponse("Failed loading Slack Webhook URL", 500)
	}

	fmt.Println("Received body: ", request.Body)

	return modules.CreateResponse("Message sent successfuly!", 200)
}

func main() {
	lambda.Start(Handler)
}
