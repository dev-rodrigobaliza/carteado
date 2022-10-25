package main

import (
	"log"

	"github.com/dev-rodrigobaliza/carteado/internal/bootstrap"
)

// these variables will be updated by the build script
const (	
	APPMADEBY = "dev-rodrigobaliza"
	APPNAME   = "carteado"
)

var (
	appDate    = "unknown"
	appVersion = "0.0.0"
	debug      = "false"
)

func main() {
	dbg := (debug == "true")

	err := bootstrap.NewApp(APPMADEBY, APPNAME, appVersion, appDate, dbg)
	if err != nil {
		log.Printf("error creating server: %s", err.Error())
	}
}
