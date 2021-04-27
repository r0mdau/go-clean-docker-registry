package main

import (
	"log"
	"os"
)

func main() {
	app := initCmdApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
