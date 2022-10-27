package main

import (
	"embed"
	"log"

	"github.com/dev-rodrigobaliza/carteado/internal/bootstrap"
)

const (
	APPMADEBY = "dev-rodrigobaliza"
	APPNAME   = "carteado"
)

var (
	//go:embed all:dist
	assets embed.FS
	// these variables will be updated by the build script
	appDate    = "unknown"
	appVersion = "0.0.0"
	debug      = "false"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func main() {
	dbg := (debug == "true")

	err := bootstrap.NewApp(APPMADEBY, APPNAME, appVersion, appDate, dbg, assets)
	if err != nil {
		log.Printf("error creating server: %s", err.Error())
	}
}
