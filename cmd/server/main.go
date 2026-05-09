package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/ajc133/calendarproxy/pkg/db"
	"github.com/ajc133/calendarproxy/pkg/handlers"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	httpInterface = os.Getenv("HTTP_INTERFACE")
	httpPort      = os.Getenv("HTTP_PORT")
)

func init() {
	if httpInterface == "" {
		httpInterface = "127.0.0.1"
	}
	if httpPort == "" {
		httpPort = "8080"
	}
}

//go:embed static
var server embed.FS

func main() {
	log.Printf("Starting application")
	db.InitDB()

	fs, err := static.EmbedFolder(server, "static")
	if err != nil {
		log.Fatalf("%s", err)
	}
	router := gin.Default()
	router.Use(static.Serve("/", fs))
	router.GET("/calendars/:id", handlers.GetCalendarByID)
	router.POST("/calendars", handlers.CreateCalendar)
	router.PATCH("/calendars", handlers.UpdateCalendar)
	socket := fmt.Sprintf("%s:%s", httpInterface, httpPort)
	router.Run(socket)
}
