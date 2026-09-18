package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"atlas/internal/core/tracking"
	"atlas/internal/shared/config"
	"atlas/internal/shared/ws"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	hub    *ws.Hub
	router *gin.Engine
}

func New(cfg *config.Config, db *gorm.DB) (*App, error) {
	router := gin.New()
	router.Use(gin.Recovery())

	a := &App{
		cfg:    cfg,
		db:     db,
		hub:    ws.NewHub(),
		router: router,
	}

	if err := a.registerModules(); err != nil {
		return nil, err
	}
	a.registerSystemRoutes()
	return a, nil
}

// Router expone el engine (util para tests con httptest).
func (a *App) Router() *gin.Engine { return a.router }

func (a *App) Run() error {
	log.Printf("atlas escuchando en :%s", a.cfg.Port)
	return a.router.Run(":" + a.cfg.Port)
}

// registerModules es el unico lugar donde se cablean los modulos de negocio.
func (a *App) registerModules() error {
	trackingModule, err := tracking.New(a.db, a.hub)
	if err != nil {
		return fmt.Errorf("modulo tracking: %w", err)
	}
	trackingModule.Register(a.router)

	// Modulos siguientes (fleet, geofencing, alerts...) se agregan aqui.
	return nil
}

func (a *App) registerSystemRoutes() {
	a.router.GET("/health", a.health)
}

func (a *App) health(c *gin.Context) {
	dbEstado := "ok"
	if sqlDB, err := a.db.DB(); err != nil || sqlDB.Ping() != nil {
		dbEstado = "error"
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"service":  "atlas",
		"database": dbEstado,
	})
}
