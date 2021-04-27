package main

import (
	"github.com/r0mdau/go-clean-docker-registry/cmd"
	"log"
	"os"
)

func main() {
	app := cmd.CreateApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
