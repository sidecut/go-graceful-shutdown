package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/10", func(c echo.Context) error {
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusOK, echo.Map{"a": "b", "c": "d"})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
