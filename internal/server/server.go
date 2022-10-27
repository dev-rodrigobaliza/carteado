package server

import (
	"embed"
	"net/http"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/handlers"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/repositories"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core/saloon"
	"github.com/dev-rodrigobaliza/carteado/internal/database"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	config *config.App
	db     *database.Database
	app    *fiber.App
	saloon *saloon.Saloon
}

func New(cfg *config.App, assets embed.FS) *Server {
	app := fiber.New(fiber.Config{
		AppName:               cfg.Name,
		ServerHeader:          cfg.Name,
		DisableStartupMessage: !cfg.Debug,
		DisableDefaultDate:    !cfg.Debug,
		IdleTimeout:           utils.StringToDuration(cfg.HTTP.IdleTimeout),
		ReadTimeout:           utils.StringToDuration(cfg.HTTP.ReadTimeout),
		WriteTimeout:          utils.StringToDuration(cfg.HTTP.WriteTimeout),
	})

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}

	s := &Server{
		config: cfg,
		app:    app,
		db:     db,
	}

	app.Use(s.timing())
	app.Use(s.versioning())
	if cfg.HTTP.Limiter.Enabled {
		app.Use(limiter.New(limiter.Config{
			Next:              s.limiterNext,
			Max:               cfg.HTTP.Limiter.MaxRequests,
			Expiration:        utils.StringToDuration(cfg.HTTP.Limiter.Expiration),
			LimiterMiddleware: limiter.SlidingWindow{},
			LimitReached:      s.limiterLimitReached,
		}))
	}
	app.Use(recover.New(recover.Config{}))

	app.Get("/health", s.getHealth)

	api := s.app.Group("api")
	appRepository := repositories.NewAppRepository(db.DB)
	appService := services.NewAppService(appRepository)
	handlers.Load(appService, api)

	ws := app.Group("ws")
	s.initSaloon(ws, appService)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(assets),
		PathPrefix:   "dist",
	}))

	app.Use(s.error404())

	return s
}

func (s *Server) Start() {
	go func() {
		err := s.app.Listen(s.config.HTTP.Address)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *Server) Stop() error {
	s.saloon.Stop()
	return s.app.Shutdown()
}
