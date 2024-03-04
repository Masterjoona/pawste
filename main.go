package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type File struct {
	ID   int
	Name string
	Size int
	Blob []byte
}

type Paste struct {
	ID             int
	PasteName      string
	Expire         string
	Privacy        string
	ReadCount      int
	ReadLast       string
	BurnAfter      int
	Content        string
	UrlRedirect    int
	Syntax         string
	HashedPassword string
	Files          []File
	CreatedAt      string
	UpdatedAt      string
}

var PasteDB *sql.DB

func main() {
	InitConfig()
	println("Using version: " + PawsteVersion)
	PasteDB = CreateOrLoadDatabase(Config.IUnderstandTheRisks)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/css", "./css")
	r.Static("/js", "./js")
	r.Static("/fonts", "./fonts")

	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/static/suzume.png", "./static/suzume.png")

	r.GET("/", HandlePage(gin.H{}, nil, ""))

	r.GET("/p/:pasteName", HandlePastePage)
	r.GET("/p", RedirectHome)

	r.GET("/u/:pasteName", Redirect)
	r.GET("/u", RedirectHome)

	r.GET("/r/:pasteName", HandleRaw)
	r.GET("/r", RedirectHome)

	r.POST("/submit", HandleSubmit)
	r.PATCH("/submit/:pasteName", HandleUpdate)

	r.GET("/guide", HandlePage(gin.H{"Guide": true}, nil, ""))
	r.GET("/admin", HandlePage(gin.H{"Admin": true}, AdminHandler, "Pastes"))
	r.GET("/list", HandlePage(gin.H{"List": true}, ListHandler, "PasteLists"))

	r.Run(Config.Port)
}
