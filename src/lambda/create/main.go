package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/datareal/repot/src/modules"
	"github.com/datareal/repot/src/modules/database"
)

type message struct {
	Title      string
	Quantity   string
	UpdateQtd  string
	AddedQtd   string
	OfflineQtd string
}

func createMessageString(Message message) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n", Message.Title, Message.Quantity, Message.UpdateQtd, Message.AddedQtd, Message.OfflineQtd)
}

func handleRequest(context context.Context, event interface{}) (modules.Response, error) {
	items, err := database.ScanAll()
	if err != nil {
		fmt.Println("Error:queryItems ", err)
		return modules.CreateResponse(modules.ResponseMessage{
			Message: "Error when creating the message",
		}, 500)
	}
	fmt.Println("Items len: ", len(items))

	var updated int = 0
	var added int = 0
	var offline int = 0
	for _, item := range items {
		if item.Status == "UPDATE" || item.Status == "REPROCESS" {
			updated++
		} else if item.Status == "OFFLINE" {
			offline++
		} else {
			added++
		}
	}

	var ResponseMessage = createMessageString(message{
		Title:      fmt.Sprintf("*Report crawler* - %s", time.Now().Format("01-02-2006")),
		Quantity:   fmt.Sprintf("*Quantidade de imóveis*: %d", len(items)),
		UpdateQtd:  fmt.Sprintf("*Quantidade atualizada*: %d", updated),
		AddedQtd:   fmt.Sprintf("*Quantidade adicionada*: %d", added),
		OfflineQtd: fmt.Sprintf("*Quantidade offline*: %d", offline),
	})

	return modules.CreateResponse(modules.ResponseMessage{
		Message: ResponseMessage,
	}, 200)
}

func main() {
	lambda.Start(handleRequest)
}
