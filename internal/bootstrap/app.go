package bootstrap

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/internal/server"
)

func NewApp(madeBy, name, version, date string, debug bool) error {
	start := time.Now()

	c := config.NewApp(name, version, date, madeBy, debug)
	err := c.LoadFromFile(consts.APP_CONFIG_FILENAME)
	if err != nil {
		return err
	}
	c.StartedAt = start

	showInfo(c)
	initSecurity(c)

	s := server.New(c)
	s.Start()
	log.Printf("*** server started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	close(sig)

	s.Stop()
	log.Printf("*** server stopped after %v", time.Since(start))

	return nil
}
