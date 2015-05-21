package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/uhuraapp/uhura-images/database"
	"willnorris.com/go/imageproxy"
)

var DB gorm.DB

func cors(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS, GET, POST, PUT")
	rw.Header().Set("Access-Control-Expose-Headers", "Content-Length")
}

func main() {
	DB = database.NewPostgresql()

	e := echo.New()
	e.Use(mw.Logger)
	e.Use(cors)

	e.Options("/*", func(*echo.Context) error {
		return nil
	})

	e.Get("/cache/:id", get)
	e.Get("/resolve", resolve)

	e.Run(":" + os.Getenv("PORT"))
}

func get(c *echo.Context) error {
	var image database.Image

	err := DB.Table("images").Where("id = ?", c.P(0)).Find(&image).Error

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	c.Response.Header().Add("Cache-Control", "public, max-age=31536000")
	c.Response.Header().Add("Last-Modified", image.UpdatedAt.Format(time.RFC1123))
	c.Response.Header().Add("Expires", "Mon, 25 Jun 2030 21:31:12 GMT")

	if ifModifiedSince := c.Request.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
		ifModifiedSinceTime, err := time.Parse(time.RFC1123, ifModifiedSince)
		updatedAt, _ := time.Parse(time.RFC1123, image.UpdatedAt.Format(time.RFC1123))
		if err == nil && ifModifiedSinceTime.Sub(updatedAt) < 1 {
			c.NoContent(304)
			return nil
		}
	}

	c.Response.Write(image.Data)
	return nil
}

func resolve(c *echo.Context) error {
	var id int64
	var ids []int64
	url := c.Request.URL.Query().Get("url")
	err := DB.Select("id").Table("images").Where("url = ?", url).Pluck("id", &ids).Error
	log.Println(ids, err)

	if len(ids) == 0 {
		requestedImage, requestErr := requestURL(url)
		if requestErr != nil {
			return c.NoContent(http.StatusNotFound)
		}

		image, err := save(url, requestedImage)

		if err != nil {
			return c.String(500, err.Error())
		}

		id = image.Id
	} else {
		id = ids[0]
	}

	idS := strconv.Itoa(int(id))

	c.Redirect(http.StatusMovedPermanently, "/cache/"+idS)
	return nil
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

func save(url string, _image []byte) (*database.Image, error) {
	image, err := imageproxy.Transform(_image, resizeOptions())

	if err != nil {
		return nil, err
	}

	model := database.Image{
		Url:  url,
		Data: image,
	}

	log.Println("Saving image")
	err = DB.Table("images").Save(&model).Error

	if err != nil {
		return nil, err
	}

	return &model, err
}

func resizeOptions() imageproxy.Options {
	return imageproxy.ParseOptions("250x")
}
