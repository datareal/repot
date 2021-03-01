package modules

import (
	"github.com/aws/aws-lambda-go/events"
)

// Request is used in our lambda functions to define APIGateway Requests
type Request events.APIGatewayProxyRequest
