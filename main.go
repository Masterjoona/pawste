package main

import (
	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/route"
	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
)

func main() {
	config.Vars.InitConfig()
	rlog.Info("Starting Pawste " + config.PawsteVersion)
	database.CreateOrLoadDatabase()

	r := gin.New()

	route.SetupMiddleware(r)
	route.SetupErrorHandlers(r)

	route.SetupPublicRoutes(r)
	route.SetupPasteRoutes(r)
	route.SetupRedirectRoutes(r)
	route.SetupEditRoutes(r)
	route.SetupAdminRoutes(r)

	r.Run(config.Vars.Port)
}
