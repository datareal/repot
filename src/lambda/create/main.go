package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/datareal/repot/src/modules"
	"github.com/datareal/repot/src/modules/database"
)

type repotTarget struct {
	Domain string `json:"real_estate"`
	Date   string `json:"date"`
}

type message struct {
	Title      string
	RealEstate string
	Quantity   string
	UpdateQtd  string
	AddedQtd   string
}

func createMessageString(Message message) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n", Message.Title, Message.RealEstate, Message.Quantity, Message.UpdateQtd, Message.AddedQtd)
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

	var requestTarget repotTarget
	err := json.Unmarshal([]byte(request.Body), &requestTarget)
	if err != nil {
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Error when reading request data",
		}, http.StatusInternalServerError)
	}

	fmt.Println("Creating message for ", requestTarget)

	items, err := database.QueryItems(requestTarget.Domain, requestTarget.Date)
	if err != nil {
		fmt.Println("Error:queryItems ", err)
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Error when creating the message",
		}, 500)
	}
	fmt.Println("Items: ", items)
	fmt.Println("Items Len: ", len(items))

	var updated int = 0
	var added int = 0

	for _, item := range items {
		if item.State == "UPDATED" {
			updated++
		} else {
			added++
		}
	}

	var ResponseMessage = createMessageString(message{
		Title:      fmt.Sprintf("*Report crawler* - %s", requestTarget.Date),
		RealEstate: fmt.Sprintf("*Imobiliária*: %s", requestTarget.Domain),
		Quantity:   fmt.Sprintf("*Quantidade de imóveis*: %d", len(items)),
		UpdateQtd:  fmt.Sprintf("*Quantidade atualizada*: %d", updated),
		AddedQtd:   fmt.Sprintf("*Quantidade adicionada*: %d", added),
	})

	return modules.CreateResponse(modules.ResponseMessage{
		Message: ResponseMessage,
	}, 200)
}

func main() {
	lambda.Start(Handler)
}
