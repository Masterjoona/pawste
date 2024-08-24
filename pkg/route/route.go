package route

import (
	_ "embed"
	"net/http"
	"regexp"
	"time"

	"github.com/Masterjoona/pawste/pkg/build"
	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/handling"
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"

	ginzap "github.com/gin-contrib/zap"
)

//go:embed favicon.ico
var favicon []byte

var wrapMiddleware = func(middleware func(http.Handler) http.Handler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			golte.AddLayout(r, "layout/main",
				map[string]any{
					"AnimeGirls": config.Vars.AnimeGirlMode,
					"PublicList": config.Vars.PublicList,
				},
			)
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

func SetupMiddleware(r *gin.Engine) {
	var golteRegex = regexp.MustCompile(`/golte_/(?:entries|chunks|assets)/.{1,60}|/favicon`)
	r.Use(ginzap.GinzapWithConfig(config.Logger.Desugar(), &ginzap.Config{
		UTC:             true,
		TimeFormat:      time.RFC3339,
		SkipPathRegexps: []*regexp.Regexp{golteRegex},
	}))
	r.Use(ginzap.RecoveryWithZap(config.Logger.Desugar(), true))
	r.Use(wrapMiddleware(build.Golte))
}

func SetupPublicRoutes(r *gin.Engine) {
	r.GET("/", handling.HandleNewPage)
	if config.Vars.PublicList {
		r.GET("/list", func(c *gin.Context) {
			golte.RenderPage(c.Writer, c.Request, "page/list", map[string]any{
				"pastes": database.GetAllPublicPastes(),
			})
		})
	}
	r.GET("/about", page("page/about"))
	r.GET("/guide", page("page/guide"))
	r.GET("/favicon", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", favicon)
	})
}

func SetupPasteRoutes(r *gin.Engine) {
	pasteGroup := r.Group("/p")
	{
		pasteGroup.GET("/:pasteName", handling.HandlePaste)
		pasteGroup.GET("/:pasteName/raw", handling.HandlePasteRaw)
		pasteGroup.GET("/:pasteName/json", handling.HandlePasteJson)
		pasteGroup.DELETE("/:pasteName", handling.HandlePasteDelete)
		pasteGroup.POST("/new", handling.HandleSubmit)
		pasteGroup.PATCH("/:pasteName", handling.HandleEditJson)
		pasteGroup.GET("/:pasteName/f/:fileName", handling.HandleFile)
		pasteGroup.GET("/:pasteName/f/:fileName/json", handling.HandleFileJson)
	}
}

func SetupRedirectRoutes(r *gin.Engine) {
	redirectGroup := r.Group("/")
	{
		redirectGroup.GET("u/:pasteName", handling.Redirect)
		redirectGroup.GET("u", handling.RedirectHome)
		redirectGroup.GET("p", handling.RedirectHome)
		redirectGroup.GET("r", handling.RedirectHome)
		redirectGroup.GET("e", handling.RedirectHome)
	}
}

func SetupEditRoutes(r *gin.Engine) {
	r.GET("/e/:pasteName", handling.HandleEdit)
}

func SetupAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("", page("page/admin"))
		adminGroup.GET("/json", handling.HandleAdminJson)
		adminGroup.POST("/reload-config", config.Vars.ReloadConfig)
	}
}

func SetupErrorHandlers(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		err := "Page not found"
		config.Logger.Error(err)
		golte.RenderPage(c.Writer, c.Request, "page/error", map[string]any{
			"error": err,
		})
	})

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Error()
			config.Logger.Error(err)
			golte.RenderPage(c.Writer, c.Request, "page/error", map[string]any{
				"error": err,
			})
			c.Abort()
		}
	})
}
