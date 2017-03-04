package main

import (
	"log"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/luizbafilho/chart-server/publisher"
	_ "github.com/luizbafilho/chart-server/storage/s3"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	publisher, err := publisher.New()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/index.yaml", func(c echo.Context) error {

		b, err := yaml.Marshal(publisher.GetIndex())
		if err != nil {
			return err
		}

		return c.Blob(http.StatusOK, "text/vnd.yaml", b)
	})

	e.PUT("/charts", func(c echo.Context) error {
		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		file := form.File["charts"][0]

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		if err := publisher.Publish(file.Filename, src); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
