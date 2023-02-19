package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	fmt.Println("yo 😮")
	// var router *gin.Engine = gin.Default()
	// router.GET("/yo", redirect)

	// router.Run()
	fmt.Println(createItem("teste"))
	item, err := readItem("boom")
	if err != nil {
		log.Fatalf("Failed to fetch item: %s", err)
	}
	fmt.Println(item)
	fmt.Println(item.ShortUrl)
	fmt.Println(item.LongUrl)
	getAllItems("url_shortener__entries")

	fmt.Println("=============================================================")

	tableName := "url_shortener__entries"

	newItem := Item{
		ShortUrl: "ohoo?",
		LongUrl: "?hooo",
	}
	err = createItem2(tableName, newItem)
	if err != nil {
		log.Fatalf("Failed to create item: %s", err)
	}

	input := map[string]types.AttributeValue{
		"short_url": &types.AttributeValueMemberS{Value: "ohoo?"},
	}
	result, err := getItem(tableName, input)
	if err != nil {
		log.Fatalf("Failed to fetch item: %s", err)
	}
	gotItem := &Item{}
	err = attributevalue.UnmarshalMap(result, gotItem)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	fmt.Println(gotItem)

	allItemsResult, err := getAllItems2(tableName)
	if err != nil {
		log.Fatalf("Failed to fetch all items from table %s: %s", tableName, err)
	}
	allItems := &[]Item{}
	err = attributevalue.UnmarshalListOfMaps(allItemsResult, allItems)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	fmt.Println(allItems)
}

// func redirect(c *gin.Context) {
// 	c.Redirect(http.StatusFound, "https://fun.pyoh.dev/")
// }

type Item struct {
	ShortUrl string `json:"short_url" dynamodbav:"short_url"`
	LongUrl  string `json:"long_url" dynamodbav:"long_url"`
}

func createItem(input string) Item {
	tableName := "url_shortener__entries"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	output := []rune(input)
	rand.Shuffle(len(output), func(i, j int) {
		output[i], output[j] = output[j], output[i]
	})
	item := Item{
		ShortUrl: input,
		LongUrl:  string(output),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}
	tableInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(context.TODO(), tableInput)
	if err != nil {
		panic(err)
	}

	return item
}

func readItem(shortUrl string) (*Item, error) {
	tableName := "url_shortener__entries"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"short_url": &types.AttributeValueMemberS{Value: shortUrl},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find item with short_url '" + shortUrl + "'"
		return nil, errors.New(msg)
	}

	item := &Item{}
	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item, nil
}

func getAllItems(tableName string) *[]Item {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		panic(err)
	}

	items := &[]Item{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println(items)
	return items
}
