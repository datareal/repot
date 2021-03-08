package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const awsRegion = "us-east-1"
const tableName = "datareal-crawler-reports"

var database = dynamodb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

// Item is the item structure from DynamoDB `datareal-crawler-reports` table
type Item struct {
	State string `json:"state"`
}

func queryItems(Date string, LastEvaluatedKey map[string]*dynamodb.AttributeValue) ([]Item, map[string]*dynamodb.AttributeValue, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName:         aws.String(tableName),
		IndexName:         aws.String("date-index"),
		ExclusiveStartKey: LastEvaluatedKey,
		KeyConditions: map[string]*dynamodb.Condition{
			"date": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(Date),
					},
				},
			},
		},
	}

	var response, err = database.Query(queryInput)
	if err != nil {
		return []Item{}, nil, err
	}
	if len(response.Items) == 0 {
		return []Item{}, nil, nil
	}
	fmt.Println("LastEvaluatedKey from query is", response.LastEvaluatedKey)

	var items []Item
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &items)
	if err != nil {
		return []Item{}, nil, err
	}

	return items, response.LastEvaluatedKey, nil
}

// QueryAll retrieves all the item_reports from the DB following the query
func QueryAll(Date string) ([]Item, error) {
	var result []Item
	var PaginationDone bool = false
	var ExclusiveStartKey map[string]*dynamodb.AttributeValue = nil

	for !PaginationDone {
		items, LastEvaluatedKey, err := queryItems(Date, ExclusiveStartKey)
		if err != nil {
			return []Item{}, err
		}
		result = append(result, items...)

		if LastEvaluatedKey == nil {
			fmt.Println("The final result lenght is", len(result))
			PaginationDone = true
		} else {
			ExclusiveStartKey = LastEvaluatedKey
		}
	}

	return result, nil
}

// ScanItems retrieves all the items from the desired table
func ScanItems() ([]Item, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := database.Scan(input)
	if err != nil {
		return []Item{}, err
	}
	if len(result.Items) == 0 {
		return []Item{}, nil
	}

	var items []Item
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return []Item{}, err
	}

	return items, nil
}
