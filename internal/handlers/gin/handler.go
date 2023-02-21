package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrofbo/url_shortener/dynamodb"
	"github.com/pedrofbo/url_shortener/internal"
)

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
		internal.Error.Println(err)
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
		internal.Error.Println(err)
		msg := fmt.Sprintf("{\"error\":\"The request body must be a json with the key "+
			"`url` and the value of the url to be shortened. %s\"}", err)
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	shortUrl, err := internal.GenerateShortUrl(requestBody.Url, config.ParticlesTableName)
	if err != nil {
		internal.Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}

	item := internal.Item{
		ShortUrl: *shortUrl,
		LongUrl:  requestBody.Url,
	}
	err = dynamodb.CreateItem(config.EntriesTableName, item)
	if err != nil {
		internal.Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	redirectUrl, err := internal.JoinUrl(config.BaseEndpoint, *shortUrl)
	if err != nil {
		internal.Error.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	response := internal.Item{
		ShortUrl: redirectUrl,
		LongUrl:  requestBody.Url,
	}
	c.IndentedJSON(http.StatusCreated, response)
}
