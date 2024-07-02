package main

import (
	"database/sql"
	"net/http"

	"github.com/Masterjoona/pawste/pkg/build"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/handling"
	"github.com/Masterjoona/pawste/pkg/shared/config"

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
	r.Use(layout("layout/main"))

	r.GET("/", page("page/index"))
	r.GET("/new", page("page/new"))
	r.GET("/list", func(ctx *gin.Context) {
		golte.RenderPage(ctx.Writer, ctx.Request, "page/list", map[string]any{
			"pasteArr":    database.GetPublicPastes(),
			"redirectArr": database.GetPublicRedirects(),
		})
	})
	r.GET("/test", func(ctx *gin.Context) {
		paste, _ := database.GetPasteByName("meerkat-meerkat-meerkat")
		golte.RenderPage(ctx.Writer, ctx.Request, "page/paste", map[string]any{
			"paste": paste,
		})
	})
	r.GET("/about", page("page/about"))
	r.GET("/guide", page("page/guide"))

	r.LoadHTMLGlob("oldweb/templates/*")

	r.Static("/css", "./oldweb/css")
	r.Static("/js", "./oldweb/js")
	r.Static("/fonts", "./oldweb/fonts")

	r.StaticFile("/favicon.ico", "./oldweb/static/favicon.ico")
	r.StaticFile("/static/suzume.png", "./oldweb/static/suzume.png")

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

	r.GET("/admin", handling.HandlePage(gin.H{"Admin": true}, handling.AdminHandler, "PasteLists"))
	r.POST("/admin/reload-config", config.Config.ReloadConfig)

	// for testing purposes
	r.GET("/old", handling.HandlePage(gin.H{}, nil, ""))

	r.Run(config.Config.Port)
}
