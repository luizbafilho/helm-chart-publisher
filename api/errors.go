package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/luizbafilho/chart-server/publisher"
)

type ErrResponse map[string]interface{}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	response := ErrResponse{}

	switch v := err.(type) {
	case publisher.ResourceNotFoundErr:
		code = http.StatusNotFound
		response = ErrResponse{"error": v.Error()}
	default:
		response = ErrResponse{"error": v.Error()}
	}

	c.JSON(code, response)
}
