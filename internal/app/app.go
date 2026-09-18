package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"atlas/internal/shared/config"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	router *gin.Engine
}

func New(cfg *config.Config, db *gorm.DB) *App {
	router := gin.New()
	router.Use(gin.Recovery())

	a := &App{cfg: cfg, db: db, router: router}
	a.registerRoutes()
	return a
}

// Router expone el engine (util para tests con httptest).
func (a *App) Router() *gin.Engine { return a.router }

func (a *App) Run() error {
	log.Printf("atlas escuchando en :%s", a.cfg.Port)
	return a.router.Run(":" + a.cfg.Port)
}

func (a *App) registerRoutes() {
	// Endpoints del sistema
	a.router.GET("/health", a.health)

	// Modulos de negocio (se agregan aqui, uno por linea):
	// tracking.New(a.db).Register(a.router)
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
