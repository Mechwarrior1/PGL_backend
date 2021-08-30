package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Mechwarrior1/PGL_backend/controller"
	"github.com/Mechwarrior1/PGL_backend/model"
	"github.com/Mechwarrior1/PGL_backend/mysql"
	"github.com/Mechwarrior1/PGL_backend/postgres"
	"github.com/Mechwarrior1/PGL_backend/word2vec"

	_ "github.com/go-sql-driver/mysql" // go mod init api_server.go
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartServer() (http.Server, *echo.Echo, *model.DBHandler, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	embed := word2vec.GetWord2Vec("word2vec/GoogleNews-vectors-negative300-SLIM.bin") // credits: https://github.com/eyaler/word2vec-slim

	err := godotenv.Load("go.env")
	if err != nil {
		fmt.Println(err)
	}

	port := os.Getenv("PORT")
	database := os.Getenv("DBTYPE")

	dbHandler1 := model.DBHandler{nil,
		os.Getenv("API_KEY"),
		false}

	switch database {
	case "postgres":
		dbHandler1.DBController = postgres.OpenDB()
	case "mysql":
		dbHandler1.DBController = mysql.OpenDB()
	}

	e := echo.New()
	e.GET("/api/v0/check", func(c echo.Context) error {
		return controller.PwCheck(c, &dbHandler1)
	})

	e.GET("/api/v0/comment/:id", func(c echo.Context) error {
		return controller.GetAllComment(c, &dbHandler1, embed)
	})

	e.GET("/api/v0/index", func(c echo.Context) error {
		return controller.GetAllListingIndex(c, &dbHandler1, embed)
	})

	e.GET("/api/v0/listing", func(c echo.Context) error {
		return controller.GetAllListing(c, &dbHandler1)
	})

	e.GET("/api/v0/username/:username", func(c echo.Context) error {
		return controller.CheckUsername(c, &dbHandler1)
	})

	e.POST("/api/v0/db/info", func(c echo.Context) error {
		return controller.GenInfoPost(c, &dbHandler1)
	})

	e.POST("/api/v0/db/signup", func(c echo.Context) error {
		return controller.Signup(c, &dbHandler1)
	})

	e.PUT("/api/v0/db/completed/:id", func(c echo.Context) error {
		return controller.Completed(c, &dbHandler1)
	})

	e.GET("/api/v0/db/info", func(c echo.Context) error {
		return controller.GenInfoGet(c, &dbHandler1)
	})

	e.PUT("/api/v0/db/info", func(c echo.Context) error {
		return controller.GenInfoPut(c, &dbHandler1)
	})

	e.GET("/api/v0/health", controller.HealthCheckLiveness)

	e.GET("/api/v0/ready", func(c echo.Context) error {
		return controller.HealthCheckReadiness(c, &dbHandler1)
	})

	// e.DELETE("/api/v0/db/info", func(c echo.Context) error {
	// 	return controller.GenInfoDelete(c, &dbHandler1)
	// })

	// go routine for checking mysql connection, will update readiness if connected
	go func(dbHandler *model.DBHandler) {
		for {
			_, err1 := dbHandler1.DBController.GetSingleRecord("ItemListing", "WHERE ID", "000001")
			if err1 != nil {
				dbHandler.ReadyForTraffic = false
				fmt.Println("unable to contact mysql server", err1.Error())
			} else {
				dbHandler.ReadyForTraffic = true
			}
			time.Sleep(120 * time.Second)
		}
	}(&dbHandler1)

	fmt.Println("listening at port " + port)
	s := http.Server{Addr: ":" + port, Handler: e}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	return s, e, &dbHandler1, nil
}

func main() {
	s, e, dbHandler1, _ := StartServer()
	defer dbHandler1.DBController.ReturnDB().Close()
	// if err := s.ListenAndServeTLS("secure//cert.pem", "secure//key.pem"); err != nil && err != http.ErrServerClosed {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}
