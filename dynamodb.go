package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// createItem creates the given item in the given DynamoDB table.
func createItem2(tableName string, item interface{}) error {
	client, err := getDynamodbClient()
	if err != nil {
		panic(err)
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("Got error marshalling new item: %s", err)
	}
	tableInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(context.TODO(), tableInput)
	if err != nil {
		return err
	}

	return nil
}

// getItem returns an item from the given DynamoDB table with the given
// input, which will be used as a key to query the table. The returned value
// is direct assigned from the output of the GetItem function from the
// DynamoDB client.
func getItem(tableName string, input map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	client, err := getDynamodbClient()
	if err != nil {
		return nil, err
	}

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       input,
	})
	if err != nil {
		return nil, fmt.Errorf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return nil, errors.New("Could not find item with the given input")
	}

	return result.Item, nil
}

// getAllItems returns all items from the given DynamoDB table. The returned
// value is direct assigned from the output of the ScanInput function from the
// DynamoDB client.
func getAllItems2(tableName string) ([]map[string]types.AttributeValue, error) {
	client, err := getDynamodbClient()
	if err != nil {
		return nil, err
	}

	result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// getDynamodbClient creates a DynamoDB client using the default AWS
// credentials.
func getDynamodbClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}
