package main

import (
	"os"
	"fmt"
	"encoding/base64"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/datareal/repot/src/modules"
)

func Handler(request modules.Request) (modules.Response, error) {
	return modules.CreateResponse("Creating message", 200)
}

func main() {
	lambda.Start(Handler)
}