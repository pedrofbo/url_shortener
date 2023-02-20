package lambda

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pedrofbo/url_shortener/internal"
)

func RedirectMain() {
	lambda.Start(redirectHandler)
}

func redirectHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := internal.LoadConfig()
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	shortUrl, ok := request.PathParameters["short_url"]
	if !ok {
		return events.APIGatewayProxyResponse{}, errors.New("Invalid value for short_url")
	}

	result, err := internal.GetLongUrl(config.EntriesTableName, shortUrl)
	if err != nil {
		internal.Error.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": result.LongUrl,
		},
	}, nil
}
