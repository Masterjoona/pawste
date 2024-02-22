package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Paste struct {
	ID             int
	PasteName      string
	Content        string
	Expire         string
	BurnAfter      int
	Privacy        string
	Syntax         string
	HashedPassword string
}

var PasteDB *sql.DB

func main() {
	InitConfig()
	PasteDB = CreateOrLoadDatabase(false)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/css", "./css")
	r.Static("/js", "./js")
	r.Static("/fonts", "./fonts")

	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/static/suzume.png", "./static/suzume.png")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})

	r.POST("/submit", HandleSubmit)
	r.GET("/p/:pasteName", HandlePastePage)
	r.GET("/list", HandleListPage)
	r.GET("/guide", HandleGuidePage)

	r.Run(":9454")
}
