package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const awsRegion = "us-east-1"
const tableName = "datareal-crawler-reports"

var database = dynamodb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

// Item is the iten structure from DynamoDB `datareal-crawler-reports` table
type Item struct {
	ID         string `json:"id"`
	RealEstate string `json:"real_estate"`
	URL        string `json:"url"`
	Date       string `json:"date"`
	State      string `json:"state"`
}

// QueryItems retrieves all the item_reports from the DB following the query
func QueryItems(Domain string, Date string) ([]Item, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("real_estate-date-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"real_estate": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(Domain),
					},
				},
			},
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
		return []Item{}, err
	}
	if len(response.Items) == 0 {
		return []Item{}, nil
	}

	var items []Item
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &items)
	if err != nil {
		return []Item{}, err
	}

	return items, nil
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
