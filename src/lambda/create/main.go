package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/datareal/repot/src/modules"
)

// Handler is used by lambda start which start our code
func Handler(request modules.Request) (modules.Response, error) {
	return modules.CreateResponse("Creating message", 200)
}

func main() {
	lambda.Start(Handler)
}
