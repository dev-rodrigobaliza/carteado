package bootstrap

import (
	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/handlers"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core/saloon"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

func initSecurity(cfg *config.App) {
	security := paseto.NewPasetoMaker(cfg.Security.AccessToken, cfg.Name)
	saloon.Security = security
	handlers.Security = security
	services.Security = security
	services.ExpireTime = utils.StringToDuration(cfg.Security.ExpireTime)
}
