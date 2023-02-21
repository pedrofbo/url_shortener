package main

import (
	"fmt"

	"github.com/pedrofbo/url_shortener/internal"
	"github.com/pedrofbo/url_shortener/internal/handlers/gin"
	"github.com/pedrofbo/url_shortener/internal/handlers/lambda"
)

func main() {
	fmt.Println("yo 😮")

	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	switch config.ApiHandler {
	case "GIN":
		gin.GinMain()
	case "LAMBDA_SHORTEN":
		lambda.ShortenMain()
	case "LAMBDA_REDIRECT":
		lambda.RedirectMain()
	default:
		panic(fmt.Errorf("Handler `%s` not found", config.ApiHandler))
	}
}
