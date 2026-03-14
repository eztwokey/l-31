package api

import (
	"context"
	"net/http"
	"time"

	"github.com/eztwokey/l3-serv/internal/config"
	"github.com/eztwokey/l3-serv/internal/logic"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

type Api struct {
	server *http.Server
	engine *gin.Engine
	logic  *logic.Logic
	logger logger.Logger
}

func New(cfg *config.Config, logic *logic.Logic, logger logger.Logger) *Api {
	gin.SetMode(cfg.Api.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	if cfg.Api.GinMode == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	server := &http.Server{
		Addr:         cfg.Api.Addr,
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Api.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Api.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Api.IdleTimeout) * time.Second,
	}

	api := new(Api)
	api.server = server
	api.engine = engine
	api.logger = logger
	api.logic = logic
	api.registerRoutes()

	return api

}

func (a *Api) Run() error {
	return a.server.ListenAndServe()
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *Api) status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome",
	})
}

func (a *Api) registerRoutes() {
	a.engine.StaticFile("/", "./web/index.html")
	a.engine.GET("/status", a.status)
	a.engine.POST("/notify", a.createNotify)
	a.engine.GET("/notify/:id", a.getNotify)
	a.engine.DELETE("/notify/:id", a.cancelNotify)
}
