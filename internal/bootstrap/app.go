package bootstrap

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/internal/server"
)

func NewApp(madeBy, name, version, date string, debug bool) error {
	start := time.Now()

	c := config.NewApp(name, version, date, madeBy, debug)
	err := c.LoadFromFile("config.json")
	if err != nil {
		return err
	}

	showInfo(c)
	initSecurity(c)

	server.New(c).Start()
	log.Printf("*** server started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	close(sig)

	log.Printf("*** server stopped after %v", time.Since(start))

	return nil
}
