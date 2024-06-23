package main

import (
	"database/sql"
	"net/http"

	"github.com/Masterjoona/pawste/build"

	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
	"github.com/romana/rlog"
)

var PasteDB *sql.DB

var wrapMiddleware = func(middleware func(http.Handler) http.Handler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
		if golte.GetRenderContext(ctx.Request) == nil {
			ctx.Abort()
		}
	}
}

func main() {
	Config.InitConfig()
	rlog.Info("Starting Pawste " + PawsteVersion)
	PasteDB = CreateOrLoadDatabase(Config.IUnderstandTheRisks)

	page := func(c string) gin.HandlerFunc {
		return gin.WrapH(golte.Page(c))
	}
	layout := func(c string) gin.HandlerFunc {
		return wrapMiddleware(golte.Layout(c))
	}

	r := gin.Default()

	r.Use(wrapMiddleware(build.Golte))
	r.Use(layout("layout/main"))

	r.GET("/contact", page("page/contact"))

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

	r.GET("/e/:pasteName", HandleEdit)
	r.GET("/e", RedirectHome)

	r.POST("/submit", HandleSubmit)
	r.PATCH("/p/:pasteName", HandleUpdate)

	r.GET("/guide", HandlePage(gin.H{"Guide": true}, nil, ""))
	r.GET("/admin", HandlePage(gin.H{"Admin": true}, AdminHandler, "PasteLists"))
	r.POST("/admin/reload-config", Config.ReloadConfig)
	r.GET("/about", HandlePage(gin.H{"About": true}, nil, ""))
	r.GET("/list", HandlePage(gin.H{"List": true}, ListHandler, "PasteLists"))

	r.Run(Config.Port)
}
