package gin

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pedrofbo/url_shortener/dynamodb"
	"github.com/pedrofbo/url_shortener/internal"
)

var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

func GinMain() {
	var router *gin.Engine = gin.Default()
	router.GET("/:shortUrl", redirect)
	router.POST("/create", createShortUrl)

	router.Run()
}

func redirect(c *gin.Context) {
	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}

	var shortUrl string = c.Param("shortUrl")
	result, err := internal.GetLongUrl(config.EntriesTableName, shortUrl)
	if err != nil {
		Error.Println(err)
		c.Redirect(http.StatusFound, config.DefaultRedirectEndpoint)
	}
	c.Redirect(http.StatusFound, result.LongUrl)
}

func createShortUrl(c *gin.Context) {
	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}

	var requestBody struct {
		Url string `json:"url"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		Error.Println(err)
		msg := fmt.Sprintf("{\"error\":\"The request body must be a json with the key "+
			"`url` and the value of the url to be shortened. %s\"}", err)
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	shortUrl, err := internal.GenerateShortUrl(requestBody.Url, config.ParticlesTableName)
	if err != nil {
		Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}

	item := internal.Item{
		ShortUrl: *shortUrl,
		LongUrl:  requestBody.Url,
	}
	err = dynamodb.CreateItem(config.EntriesTableName, item)
	if err != nil {
		Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	response := internal.Item{
		ShortUrl: filepath.Join(config.BaseEndpoint + *shortUrl),
		LongUrl:  requestBody.Url,
	}
	c.IndentedJSON(http.StatusCreated, response)
}
