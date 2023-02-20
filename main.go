package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrofbo/url_shortener/dynamodb"
)

var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

func main() {
	fmt.Println("yo 😮")
	var router *gin.Engine = gin.Default()
	router.GET("/:shortUrl", redirect)
	router.POST("/create", createShortUrl)

	router.Run()
}

func redirect(c *gin.Context) {
	var shortUrl string = c.Param("shortUrl")
	result, err := GetLongUrl("url_shortener__entries", shortUrl)
	if err != nil {
		Error.Println(err)
		c.Redirect(http.StatusFound, "https://fun.pyoh.dev/")
	}
	c.Redirect(http.StatusFound, result.LongUrl)
}

func createShortUrl(c *gin.Context) {
	var requestBody struct {
		Url string `json:"url"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		Error.Println(err)
		msg := fmt.Sprintf("{\"error\":\"The request body must be a json with the key " +
			"`url` and the value of the url to be shortened. %s\"}", err)
		c.JSON(http.StatusBadRequest, msg)
		return
	}
	shortUrl, err := GenerateShortUrl(requestBody.Url, "url_shortener__particles")
	if err != nil {
		Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	item := Item{
		ShortUrl: *shortUrl,
		LongUrl: requestBody.Url,
	}
	err = dynamodb.CreateItem("url_shortener__entries", item)
	if err != nil {
		Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	response := Item{
		ShortUrl: "http://localhost:8080/" + *shortUrl,
		LongUrl: requestBody.Url,
	}
	c.IndentedJSON(http.StatusCreated, response)
}
