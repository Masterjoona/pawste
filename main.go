package main

import (
	"database/sql"
	"net/http"

	"github.com/Masterjoona/pawste/build"
	"github.com/Masterjoona/pawste/database"
	"github.com/Masterjoona/pawste/handling"
	"github.com/Masterjoona/pawste/shared/config"

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
	config.Config.InitConfig()
	rlog.Info("Starting Pawste " + config.PawsteVersion)
	PasteDB = database.CreateOrLoadDatabase(config.Config.IUnderstandTheRisks)

	page := func(c string) gin.HandlerFunc {
		return gin.WrapH(golte.Page(c))
	}
	layout := func(c string) gin.HandlerFunc {
		return wrapMiddleware(golte.Layout(c))
	}

	r := gin.Default()

	r.Use(wrapMiddleware(build.Golte))
	r.Use(layout("layout/footer"))

	r.GET("/new", page("page/new"))

	r.LoadHTMLGlob("oldweb/templates/*")

	r.Static("/css", "./oldweb/css")
	r.Static("/js", "./oldweb/js")
	r.Static("/fonts", "./oldweb/fonts")

	r.StaticFile("/favicon.ico", "./oldweb/static/favicon.ico")
	r.StaticFile("/static/suzume.png", "./oldweb/static/suzume.png")

	r.GET("/", handling.HandlePage(gin.H{}, nil, ""))

	r.GET("/p/:pasteName", handling.HandlePastePage)
	r.GET("/p/:pasteName/json", handling.HandlePasteJSON)
	r.GET("/p", handling.RedirectHome)

	r.GET("/u/:pasteName", handling.Redirect)
	r.GET("/u", handling.RedirectHome)

	r.GET("/r/:pasteName", handling.HandleRaw)
	r.GET("/r", handling.RedirectHome)

	r.GET("/e/:pasteName", handling.HandleEdit)
	r.GET("/e", handling.RedirectHome)

	r.POST("/submit", handling.HandleSubmit)
	r.PATCH("/p/:pasteName", handling.HandleUpdate)

	r.GET("/guide", handling.HandlePage(gin.H{"Guide": true}, nil, ""))
	r.GET("/admin", handling.HandlePage(gin.H{"Admin": true}, handling.AdminHandler, "PasteLists"))
	r.POST("/admin/reload-config", config.Config.ReloadConfig)
	r.GET("/about", handling.HandlePage(gin.H{"About": true}, nil, ""))
	r.GET("/list", handling.HandlePage(gin.H{"List": true}, handling.ListHandler, "PasteLists"))

	r.Run(config.Config.Port)
}
