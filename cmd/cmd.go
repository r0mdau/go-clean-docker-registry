package cmd

import (
	"fmt"
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
		Usage: "Dryrun only print future delete actions",
	}
	insecureFlag := &cli.BoolFlag{
		Name:  "insecure",
		Usage: "Disable TLS cert verification",
	}

	app.Commands = []*cli.Command{
		{
			Name:   "showimages",
			Usage:  "Show all images from your registry",
			Action: printRepositoriesList,
			Flags: []cli.Flag{
				urlFlag,
				insecureFlag,
			},
		}, {
			Name:   "showtags",
			Usage:  "Show all tags for your image",
			Action: printImageTagsList,
			Flags: []cli.Flag{
				urlFlag,
				imageFlag,
				insecureFlag,
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete all specified tags for your image",
			Action: deleteImage,
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

func verifyRegistryVersion(registry registry.Registry) {
	err := registry.VersionCheck()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func printRepositoriesList(c *cli.Context) error {
	registry := configureRegistry(c)
	verifyRegistryVersion(registry)

	repositories, err := registry.ListRepositories()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(repositories.GetRepository().List)
	return nil
}

func printImageTagsList(c *cli.Context) error {
	registry := configureRegistry(c)
	verifyRegistryVersion(registry)

	registryResponse, err := registry.ListImageTags(c.String("image"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(registryResponse.Body))
	fmt.Println("Total of", len(registryResponse.GetImage().Tags), "tags.")
	return nil
}

func deleteImage(c *cli.Context) error {
	registry := configureRegistry(c)
	verifyRegistryVersion(registry)

	cliImage := c.String("image")
	cliTag := c.String("tag")
	dryrun := c.Bool("dryrun")
	keep := c.Int("keep")

	registryResponse, err := registry.ListImageTags(cliImage)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tagsToDelete := registryResponse.GetImage().Tags
	if cliTag != "" {
		tagsToDelete, err = filter.MatchAndSortImageTags(registryResponse.GetImage().Tags, cliTag)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		tagsToDelete = tagsToDelete[:len(tagsToDelete)-keep]
	}

	if dryrun {
		fmt.Println("Dryrun, it should delete image : \""+cliImage+"\" with", len(tagsToDelete), "tags :", tagsToDelete)
		return nil
	}

	for _, tagToDelete := range tagsToDelete {
		digest, errGet := registry.GetDigestFromManifest(cliImage, tagToDelete)
		if errGet != nil {
			fmt.Println(errGet.Error())
		}
		errDel := registry.DeleteImage(cliImage, tagToDelete, digest)
		if errDel != nil {
			fmt.Println(errDel.Error())
		}
	}
	fmt.Println("Total of", len(tagsToDelete), "tags deleted.")
	return nil
}
