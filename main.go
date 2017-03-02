package main

import (
	"net/http"

	"github.com/ghodss/yaml"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var chartsIndex *repo.IndexFile

func main() {
	chartsIndex := repo.NewIndexFile()
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/index.yaml", func(c echo.Context) error {

		b, err := yaml.Marshal(*chartsIndex)
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
		files := form.File["charts"]

		index := repo.NewIndexFile()
		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			c, err := chartutil.LoadArchive(src)
			if err != nil {
				return err
			}

			hash, err := provenance.Digest(src)
			if err != nil {
				return err
			}

			fname := file.Filename

			index.Add(c.Metadata, fname, "http://charts.videos.globoi.com", hash)
		}

		chartsIndex.Merge(index)

		chartsIndex.SortEntries()

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
