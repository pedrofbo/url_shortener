package main

import (
	"fmt"

	"github.com/pedrofbo/url_shortener/internal"
	"github.com/pedrofbo/url_shortener/internal/handlers/gin"
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
	default:
		panic(fmt.Errorf("Handler `%s` not found", config.ApiHandler))
	}
}
