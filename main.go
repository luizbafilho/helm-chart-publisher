package main

import (
	"errors"
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
	e.GET("/:repo/index.yaml", func(c echo.Context) error {
		repo := c.Param("repo")

		index, err := publisher.GetIndex(repo)
		if err != nil {
			return err
		}

		b, err := yaml.Marshal(index)
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

		repos, ok := form.Value["repo"]
		if !ok {
			return errors.New("no repo provided")
		}

		if err := publisher.Publish(repos[0], file.Filename, src); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
