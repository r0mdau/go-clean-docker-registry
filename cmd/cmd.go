package cmd

import (
	"fmt"
	"github.com/korovkin/limiter"
	"github.com/r0mdau/go-clean-docker-registry/internal/filter"
	"github.com/r0mdau/go-clean-docker-registry/pkg/registry"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "go-clean-docker-registry"
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
		Usage: "Dry CreateApp only print future actions",
	}
	insecureFlag := &cli.BoolFlag{
		Name:  "insecure",
		Usage: "Disable TLS cert verification",
	}

	app.Commands = []*cli.Command{
		{
			Name:   "showimages",
			Usage:  "Show all images from your registry",
			Action: printRegistryCatalog,
			Flags: []cli.Flag{
				urlFlag,
				insecureFlag,
			},
		}, {
			Name:   "showtags",
			Usage:  "Show all tags for your image",
			Action: printRegistryTags,
			Flags: []cli.Flag{
				urlFlag,
				imageFlag,
				insecureFlag,
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete all specified tags for your image",
			Action: deleteRegistryTags,
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

func configureRegistry(c *cli.Context) registry.Registry {
	registry := registry.Registry{}
	registry.Configure(c.String("url"), c.Bool("insecure"))
	return registry
}

func printRegistryCatalog(c *cli.Context) error {
	registry := configureRegistry(c)
	catalog, err := registry.GetCatalog()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(catalog))
	return nil
}

func printRegistryTags(c *cli.Context) error {
	registry := configureRegistry(c)
	registryResponse, err := registry.GetTagsList(c.String("image"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(registryResponse.Body))
	fmt.Println("Total of", len(registryResponse.GetImage().Tags), "tags.")
	return nil
}

func deleteRegistryTags(c *cli.Context) error {
	registry := configureRegistry(c)
	imageName := c.String("image")
	tag := c.String("tag")
	dryrun := c.Bool("dryrun")

	registryResponse, err := registry.GetTagsList(imageName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var imageTagsToDelete []string
	if tag != "" {
		imageTagsToDelete, err = filter.MatchAndSortImageTags(registryResponse.GetImage().Tags, tag)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		imageTagsToDelete = imageTagsToDelete[:len(imageTagsToDelete)-c.Int("keep")]
	} else {
		imageTagsToDelete = registryResponse.GetImage().Tags
	}

	if dryrun {
		fmt.Println("Dryrun, it should delete image : \""+imageName+"\" with tags :", imageTagsToDelete)
	} else {
		limit := limiter.NewConcurrencyLimiter(10)
		for _, imageTagToDelete := range imageTagsToDelete {
			limit.Execute(func() {
				err := registry.DeleteImageTag(imageName, imageTagToDelete)
				if err != nil {
					fmt.Println(err.Error())
				}
			})
		}
		limit.Wait()
	}
	fmt.Println("Total of", len(imageTagsToDelete), "tags deleted.")
	return nil
}
