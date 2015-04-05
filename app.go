package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/uhuraapp/uhura-images/database"
	"willnorris.com/go/imageproxy"
)

var DB gorm.DB

func main() {
	DB = database.NewPostgresql()

	e := echo.New()
	e.Use(mw.Logger)

	e.Post("/cache", createImage)
	e.Get("/cache/:id", getImage)

	e.Run(":" + os.Getenv("PORT"))
}

func getImage(c *echo.Context) {
	var image database.Image

	DB.Table("images").Where("id = ?", c.P(0)).Find(&image)

	c.Response.Write(image.Data)
}

func createImage(c *echo.Context) {
	imageURL := c.Request.FormValue("url")
	if imageURL == "" {
		// error
	}

	originalImage, err := requestURL(imageURL)
	if err != nil {
		// notifyWrongImage
	}

	image, err := imageproxy.Transform(originalImage, resizeOptions())
	if err != nil {
		log.Println(err)

	}

	imageSaved := database.Image{
		Url:  imageURL,
		Data: image,
	}

	DB.Table("images").FirstOrCreate(&imageSaved)

	c.JSON(200, imageSaved.Id)
}

func requestURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Status is " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)

	return body, err
}

func resizeOptions() imageproxy.Options {
	return imageproxy.ParseOptions("250x")
}
