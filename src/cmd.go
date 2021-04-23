package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"time"
)

func initCmdApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Go Clean Docker Registry"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Value:    "http://localhost:5000",
			Usage:    "Registry url",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "image",
			Aliases:  []string{"i"},
			Usage:    "Image name to delete ie r0mdau/nodejs",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Image version tag to delete, regex possible ie \"master-.*\"",
		},
		&cli.IntFlag{
			Name:    "keep",
			Aliases: []string{"k"},
			Value: 0,
			Usage:   "Number of tags to keep, to combine with -t",
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Dry initCmdApp only print future actions",
		},
		&cli.BoolFlag{
			Name:  "insecure",
			Usage: "Disable TLS cert verification",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "show",
			Usage:   "Show all tags for your image",
			Action: func(c *cli.Context) error {
				registry := Registry{}
				registry.configure(c.String("url"), c.Bool("insecure"))
				registryResponse := registry.getApi("/v2/" + c.String("image") + "/tags/list")

				fmt.Println(string(registryResponse.Body))
				fmt.Println("Total of", len(registryResponse.getRegistryImage().Tags), "tags.")
				return nil
			},
		},
		{
			Name:    "delete",
			Usage:   "Delete all specified tags for your image",
			Action: func(c *cli.Context) error {
				//todo
				return nil
			},
		},
	}

	return app
}
