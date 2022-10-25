package bootstrap

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
)

func showInfo(cfg *config.App) {
	log.Printf("******************************************")
	log.Printf("*** %s v%s %s", strings.ToUpper(cfg.Name), cfg.Version, cfg.Date)
	log.Printf("*** © %d %s", time.Now().Year(), cfg.MadeBy)
	if cfg.Debug {
		log.Printf("*** compiled with %s", runtime.Version())
	}
	log.Printf("******************************************")
	log.Printf("*** running in OS %s", runtime.GOOS)
	log.Printf("*** running with %d processors", runtime.NumCPU())
	log.Printf("******************************************")
}
