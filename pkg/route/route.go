package route

import (
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/Masterjoona/pawste/pkg/build"
	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/handling"
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
	"github.com/romana/rlog"
)

// go:embed favicon.ico
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
	var skipGolteFiles = []string{}
	fs.WalkDir(build.Fsys, "client/golte_", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			rlog.Error(err)
			return err
		}
		if d.IsDir() {
			return nil
		}
		skipGolteFiles = append(skipGolteFiles, path[6:])
		return nil
	})
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipGolteFiles,
	}))
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
		fmt.Println(favicon)
		c.Data(http.StatusOK, "image/x-icon", favicon)
	})
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
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
		pasteGroup.POST("/:pasteName/f/:fileName", handling.HandleFilePost)
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
		rlog.Error(err)
		golte.RenderPage(c.Writer, c.Request, "page/error", map[string]any{
			"error": err,
		})
	})

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Error()
			rlog.Error(err)
			golte.RenderPage(c.Writer, c.Request, "page/error", map[string]any{
				"error": err,
			})
			c.Abort()
		}
	})
}
