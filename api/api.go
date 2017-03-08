package api

import (
	"net/http"

	"github.com/ghodss/yaml"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/luizbafilho/chart-server/publisher"
)

type API struct {
	echo      *echo.Echo
	publisher *publisher.Publisher
}

func New(publisher *publisher.Publisher) *API {
	e := echo.New()
	e.HTTPErrorHandler = CustomHTTPErrorHandler

	e.Use(middleware.Logger())

	api := &API{
		echo:      e,
		publisher: publisher,
	}

	api.registerHandlers()

	return api
}

func (api *API) Serve(address string) {
	api.echo.Logger.Fatal(api.echo.Start(address))
}

func (api *API) registerHandlers() {
	api.echo.GET("/:repo/index.yaml", api.getIndexHandler)
	api.echo.PUT("/charts", api.publishChartHandler)
}

func (api *API) getIndexHandler(c echo.Context) error {
	repo := c.Param("repo")

	index, err := api.publisher.GetIndex(repo)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(index)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "text/vnd.yaml", b)
}

func (api *API) publishChartHandler(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	if len(form.File["chart"]) < 1 {
		return BadRequestErr("no chart provided")
	}

	file := form.File["chart"][0]
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	repos, ok := form.Value["repo"]
	if !ok {
		return BadRequestErr("no repository provided")
	}

	if err := api.publisher.Publish(repos[0], file.Filename, src); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
