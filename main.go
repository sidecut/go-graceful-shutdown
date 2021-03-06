package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})
	e.GET("/delay/:delay", func(c echo.Context) error {
		delay, err := strconv.Atoi(c.Param("delay"))
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(delay) * time.Second)
		return c.JSON(http.StatusOK, echo.Map{"a": "b", "c": "d"})
	})

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	ctxShutdownTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutdownTimeout); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println(ctx.Err()) // prints "context canceled"
	stop()
}
