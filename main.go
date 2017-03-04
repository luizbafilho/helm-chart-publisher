package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/luizbafilho/chart-server/storage"
	_ "github.com/luizbafilho/chart-server/storage/s3"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//APP app name
const APP = "helm-charts-publisher"

var chartsIndex *repo.IndexFile

func init() {
	initViper()
}

func main() {
	chartsIndex := repo.NewIndexFile()

	storageType, config := getStorageAndConfig()

	store, err := storage.New(storageType, config)

	if err != nil {
		log.Fatal(err)
	}

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

			content, err := ioutil.ReadAll(src)
			if err != nil {
				return err
			}

			if err := store.Put(file.Filename, content); err != nil {
				fmt.Println("UploadFile", err)
				return err
			}

			chart, err := chartutil.LoadArchive(bytes.NewBuffer(content))
			if err != nil {
				fmt.Println("LoadArchive", err)
				return err
			}

			hash, err := provenance.Digest(src)
			if err != nil {
				fmt.Println("Digest", err)
				return err
			}

			fname := file.Filename

			index.Add(chart.Metadata, fname, "http://charts.videos.globoi.com", hash)
		}

		chartsIndex.Merge(index)

		chartsIndex.SortEntries()

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
