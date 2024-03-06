package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
)

var PasteDB *sql.DB

func main() {
	InitConfig()
	rlog.Info("Starting Pawste " + PawsteVersion)
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
	r.GET("/p/:pasteName/json", HandlePasteJSON)
	r.GET("/p", RedirectHome)

	r.GET("/u/:pasteName", Redirect)
	r.GET("/u", RedirectHome)

	r.GET("/r/:pasteName", HandleRaw)
	r.GET("/r", RedirectHome)

	r.POST("/submit", HandleSubmit)
	r.PATCH("/submit/:pasteName", HandleUpdate)

	r.GET("/guide", HandlePage(gin.H{"Guide": true}, nil, ""))
	r.GET("/admin", HandlePage(gin.H{"Admin": true}, AdminHandler, "PasteLists"))
	r.POST("/admin/reload-config", ReloadConfig)
	r.GET("/about", HandlePage(gin.H{"About": true}, nil, ""))
	r.GET("/list", HandlePage(gin.H{"List": true}, ListHandler, "PasteLists"))

	r.Run(Config.Port)
}
