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

	e.Post("/cache", create)
	e.Get("/cache/:id", get)
	e.Get("/resolve", resolve)

	e.Run(":" + os.Getenv("PORT"))
}

func get(c *echo.Context) {
	var image database.Image

	DB.Table("images").Where("id = ?", c.P(0)).Find(&image)

	c.Response.Write(image.Data)
}

func create(c *echo.Context) {
	imageURL := c.Request.FormValue("url")
	if imageURL == "" {
		// error
	}

	originalImage, err := requestURL(imageURL)
	if err != nil {
		// notifyWrongImage
	}

	imageSaved := save(imageURL, originalImage)

	c.JSON(200, imageSaved.Id)
}

func resolve(c *echo.Context) {
	var image database.Image

	url := c.Request.URL.Query().Get("url")

	err := DB.Table("images").Where("url = ?", url).Find(&image).Error
	if err != nil {
		originalImage, err2 := requestURL(url)
		if err2 != nil {
			// notifyWrongImage
		}
		save(url, originalImage)

		c.Response.Write(originalImage)
		return
	}

	c.Response.Write(image.Data)
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

func save(url string, image []byte) database.Image {
	newimage, err := imageproxy.Transform(image, resizeOptions())

	if err != nil {
		log.Println(err)
	}

	imageSaved := database.Image{
		Url:  url,
		Data: newimage,
	}

	log.Println("Saving image")
	log.Println(DB.Table("images").Where("url = ?", url).FirstOrCreate(&imageSaved).Error)

	return imageSaved
}
