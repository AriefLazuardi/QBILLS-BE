package main

import (
	"log"
	"net/http"
	"os"
	"qbills/drivers"
	"qbills/routes"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	myApp := echo.New()
	validate := validator.New()

	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to fetch .env file")
		}
	}

	drivers.ConnectDB()
	drivers.Migrate()

	myApp.GET("/home", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Q Bills API Services")
	})

	routes.AdminRoutes(myApp, drivers.DB, validate)
	routes.CashierRoutes(myApp, drivers.DB, validate)
	routes.SuperAdminRoutes(myApp, drivers.DB, validate)
	routes.ProductTypeRoutes(myApp, drivers.DB, validate)

	myApp.Pre(middleware.RemoveTrailingSlash())
	myApp.Use(middleware.CORS())
	myApp.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		},
	))


	myApp.Logger.Fatal(myApp.Start(":443"))
}
