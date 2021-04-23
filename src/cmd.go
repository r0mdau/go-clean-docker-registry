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

	urlFlag := &cli.StringFlag{
		Name:     "url",
		Aliases:  []string{"u"},
		Value:    "http://localhost:5000",
		Usage:    "Registry url",
		Required: true,
	}
	imageFlag := &cli.StringFlag{
		Name:     "image",
		Aliases:  []string{"i"},
		Usage:    "Image name to delete ie r0mdau/nodejs",
		Required: true,
	}
	tagFlag := &cli.StringFlag{
		Name:    "tag",
		Aliases: []string{"t"},
		Usage:   "Image version tag to delete, regex possible ie \"master-.*\"",
	}
	keepFlag := &cli.IntFlag{
		Name:    "keep",
		Aliases: []string{"k"},
		Value:   0,
		Usage:   "Number of tags to keep, to combine with -t",
	}
	dryrunFlag := &cli.BoolFlag{
		Name:  "dryrun",
		Usage: "Dry initCmdApp only print future actions",
	}
	insecureFlag := &cli.BoolFlag{
		Name:  "insecure",
		Usage: "Disable TLS cert verification",
	}

	app.Commands = []*cli.Command{
		{
			Name:  "showimages",
			Usage: "Show all images from your registry",
			Action: func(c *cli.Context) error {
				registry := Registry{}
				registry.configure(c.String("url"), c.Bool("insecure"))
				catalog := registry.getCatalog("/v2/_catalog")

				fmt.Println(string(catalog))
				return nil
			},
			Flags: []cli.Flag{
				urlFlag,
				insecureFlag,
			},
		}, {
			Name:  "showtags",
			Usage: "Show all tags for your image",
			Action: func(c *cli.Context) error {
				registry := Registry{}
				registry.configure(c.String("url"), c.Bool("insecure"))
				registryResponse := registry.getTagsList("/v2/" + c.String("image") + "/tags/list")

				fmt.Println(string(registryResponse.Body))
				fmt.Println("Total of", len(registryResponse.getRegistryImage().Tags), "tags.")
				return nil
			},
			Flags: []cli.Flag{
				urlFlag,
				imageFlag,
				insecureFlag,
			},
		},
		{
			Name:  "delete",
			Usage: "Delete all specified tags for your image",
			Action: func(c *cli.Context) error {
				//todo
				return nil
			},
			Flags: []cli.Flag{
				urlFlag,
				imageFlag,
				tagFlag,
				keepFlag,
				dryrunFlag,
				insecureFlag,
			},
		},
	}

	return app
}
