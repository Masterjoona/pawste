package main

import (
	"net/http"

	"github.com/Masterjoona/pawste/pkg/build"
	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/handling"

	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
	"github.com/romana/rlog"
)

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
	database.CreateOrLoadDatabase()

	page := func(c string) gin.HandlerFunc {
		return gin.WrapH(golte.Page(c))
	}
	layout := func(c string) gin.HandlerFunc {
		return wrapMiddleware(golte.Layout(c))
	}

	r := gin.Default()

	r.Use(wrapMiddleware(build.Golte))
	r.Use(layout("layout/main"))

	r.GET("/", page("page/new"))
	r.GET("/list", func(ctx *gin.Context) {
		golte.RenderPage(ctx.Writer, ctx.Request, "page/list", map[string]any{
			"pastes": database.GetAllPublicPastes(),
		})
	})
	r.GET("/about", page("page/about"))
	r.GET("/guide", page("page/guide"))

	r.GET("/p/:pasteName", handling.HandlePastePage)
	r.POST("/p/:pasteName", handling.HandlePastePage) // for auth
	r.GET("/p/:pasteName/auth", func(ctx *gin.Context) {
		golte.RenderPage(ctx.Writer, ctx.Request, "page/auth", map[string]any{
			"pasteRedir": ctx.Param("pasteName"),
		})
	})
	r.POST("/p/:pasteName/auth", handling.HandlePastePostAuth)
	r.GET("/p/:pasteName/raw", handling.HandlePasteRaw)
	r.GET("/p/:pasteName/json", handling.HandlePasteJSON)
	r.DELETE("/p/:pasteName", handling.HandlePasteDelete)
	r.POST("/p", handling.HandleSubmit)

	r.GET("/u/:pasteName", handling.Redirect)

	r.GET("/e/:pasteName", handling.HandleEdit)

	r.PATCH("/p/:pasteName", handling.HandleUpdate)

	r.POST("/admin/reload-config", config.Config.ReloadConfig)

	r.GET("/p/:pasteName/f/:fileName", handling.HandleFile)
	r.GET("/p/:pasteName/f/:fileName/json", handling.HandleFileJson)

	r.GET("/p", handling.RedirectHome)
	r.GET("/u", handling.RedirectHome)
	r.GET("/r", handling.RedirectHome)
	r.GET("/e", handling.RedirectHome)

	r.Run(config.Config.Port)
}
