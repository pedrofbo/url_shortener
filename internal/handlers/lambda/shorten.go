package lambda

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pedrofbo/url_shortener/dynamodb"
	"github.com/pedrofbo/url_shortener/internal"
)

func ShortenMain() {
	lambda.Start(shortenHandler)
}

func shortenHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := internal.LoadConfig()
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}
	var requestBody struct {
		Url string `json:"url"`
	}
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	shortUrl, err := internal.GenerateShortUrl(requestBody.Url, config.ParticlesTableName)
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	item := internal.Item{
		ShortUrl: *shortUrl,
		LongUrl:  requestBody.Url,
	}
	err = dynamodb.CreateItem(config.EntriesTableName, item)
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}
	response := internal.Item{
		ShortUrl: filepath.Join(config.BaseEndpoint, *shortUrl),
		LongUrl:  requestBody.Url,
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body: string(responseBody),
	}, nil
}
