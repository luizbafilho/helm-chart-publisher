package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/luizbafilho/chart-server/publisher"
)

type BadRequestErr string

func (e BadRequestErr) Error() string {
	return string(e)
}

type ErrResponse map[string]interface{}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	response := ErrResponse{}

	switch err.(type) {
	case publisher.ResourceNotFoundErr:
		code = http.StatusNotFound
	case BadRequestErr:
		code = http.StatusBadRequest
	}

	response = ErrResponse{"error": err.Error()}

	c.JSON(code, response)
}
