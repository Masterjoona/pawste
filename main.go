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

func page(c string) gin.HandlerFunc {
	return gin.WrapH(golte.Page(c))
}

func layout(c string) gin.HandlerFunc {
	return wrapMiddleware(golte.Layout(c))
}

func setupMiddleware(r *gin.Engine) {
	r.Use(wrapMiddleware(build.Golte))
	r.Use(layout("layout/main"))
}

func setupPublicRoutes(r *gin.Engine) {
	r.GET("/", page("page/new"))
	r.GET("/list", handling.HandleList)
	r.GET("/about", page("page/about"))
	r.GET("/guide", page("page/guide"))
}

func setupPasteRoutes(r *gin.Engine) {
	pasteGroup := r.Group("/p")
	{
		pasteGroup.GET("/:pasteName", handling.HandlePastePage)
		pasteGroup.GET("/:pasteName/raw", handling.HandlePasteRaw)
		pasteGroup.GET("/:pasteName/json", handling.HandlePasteJSON)
		pasteGroup.DELETE("/:pasteName", handling.HandlePasteDelete)
		pasteGroup.POST("/", handling.HandleSubmit)
		pasteGroup.PATCH("/:pasteName", handling.HandleUpdate)
		pasteGroup.GET("/:pasteName/f/:fileName", handling.HandleFile)
		pasteGroup.GET("/:pasteName/f/:fileName/json", handling.HandleFileJson)
	}
}

func setupRedirectRoutes(r *gin.Engine) {
	redirectGroup := r.Group("/")
	{
		redirectGroup.GET("u/:pasteName", handling.Redirect)
		redirectGroup.GET("u", handling.RedirectHome)
		redirectGroup.GET("p", handling.RedirectHome)
		redirectGroup.GET("r", handling.RedirectHome)
		redirectGroup.GET("e", handling.RedirectHome)
	}
}

func setupEditRoutes(r *gin.Engine) {
	r.GET("/e/:pasteName", handling.HandleEdit)
}
func setupAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("", handling.HandleAdmin)
		adminGroup.GET("/json", handling.HandleAdminJSON)
		adminGroup.POST("/reload-config", config.Config.ReloadConfig)
	}
}

func main() {
	config.Config.InitConfig()
	rlog.Info("Starting Pawste " + config.PawsteVersion)
	database.CreateOrLoadDatabase()

	r := gin.Default()

	setupMiddleware(r)

	setupPublicRoutes(r)
	setupPasteRoutes(r)
	setupRedirectRoutes(r)
	setupEditRoutes(r)
	setupAdminRoutes(r)

	r.Run(config.Config.Port)
}
