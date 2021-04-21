package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Clean Docker Registry"
	app.Version = "0.1.0"
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Value:   "http://localhost:5000",
			Usage:   "Registry url",
		},
		&cli.StringFlag{
			Name:    "image",
			Aliases: []string{"i"},
			Usage:   "Image name to delete ie r0mdau/nodejs",
		},
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Image version tag to delete, regex possible ie \"master-.*\"",
		},
		&cli.StringFlag{
			Name:    "keep",
			Aliases: []string{"k"},
			Usage:   "Delete older version tags than this version (semver compatible), to combine with a regex in -t",
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Dry run only print future actions",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(
			c.String("url"),
			c.String("image"),
			c.String("tag"),
			c.String("keep"),
			c.String("dryrun"),
		)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
