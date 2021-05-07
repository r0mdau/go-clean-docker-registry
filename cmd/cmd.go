package cmd

import (
	"fmt"
	"github.com/r0mdau/go-clean-docker-registry/internal/filter"
	"github.com/r0mdau/go-clean-docker-registry/pkg/registry"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

const workers = 10

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
		Usage:   "Image version tag to delete, regex possible ie \"master-*\", priority for semver",
	}
	keepFlag := &cli.IntFlag{
		Name:    "keep",
		Aliases: []string{"k"},
		Value:   0,
		Usage:   "Number of tags to keep, to combine with -t",
	}
	numberFlag := &cli.IntFlag{
		Name:  "n",
		Value: 5000,
		Usage: "Number of images to retrieve",
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
				numberFlag,
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

func verifyRegistryVersion(registry registry.Registry) {
	err := registry.VersionCheck()
	exit(err)
}

func printRepositoriesList(c *cli.Context) error {
	registry := registry.NewRegistry(c.String("url"), c.Bool("insecure"))
	verifyRegistryVersion(registry)

	repositories, err := registry.ListRepositories(c.Int("n"))
	exit(err)

	fmt.Printf(string(repositories.Body))
	fmt.Println("Total of", len(repositories.GetRepository().List), "repositories.")
	return nil
}

func printImageTagsList(c *cli.Context) error {
	registry := registry.NewRegistry(c.String("url"), c.Bool("insecure"))
	verifyRegistryVersion(registry)

	imageTags, err := registry.ListImageTags(c.String("image"))
	exit(err)

	fmt.Println(string(imageTags.Body))
	fmt.Println("Total of", len(imageTags.GetImage().Tags), "tags.")
	return nil
}

func deleteImage(c *cli.Context) error {
	registry := registry.NewRegistry(c.String("url"), c.Bool("insecure"))
	verifyRegistryVersion(registry)

	cliImage := c.String("image")
	cliTag := c.String("tag")
	dryrun := c.Bool("dryrun")
	keep := c.Int("keep")

	registryResponse, err := registry.ListImageTags(cliImage)
	exit(err)

	tagsToDelete := registryResponse.GetImage().Tags
	if cliTag != "" {
		tagsToDelete, err = filter.MatchAndSortImageTags(registryResponse.GetImage().Tags, cliTag)
		exit(err)
		tagsToDelete = tagsToDelete[:len(tagsToDelete)-keep]
	}

	if dryrun {
		fmt.Println("Dryrun, it should delete image : \""+cliImage+"\" with", len(tagsToDelete), "tags :", tagsToDelete)
		return nil
	}

	numJobs := len(tagsToDelete)
	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)

	for w := 0; w < workers; w++ {
		go wDelete(registry, cliImage, jobs, results)
	}
	for _, tagToDelete := range tagsToDelete {
		jobs <- tagToDelete
	}
	close(jobs)
	for a := 0; a < numJobs; a++ {
		<-results
	}
	fmt.Println("Total of", len(tagsToDelete), "tags deleted.")
	return nil
}

func wDelete(registry registry.Registry, image string, jobs <-chan string, results chan<- string) {
	for tag := range jobs {
		digest, errGet := registry.GetDigestFromManifest(image, tag)
		if errGet != nil {
			fmt.Println(errGet.Error())
			results <- tag
			continue
		}
		fmt.Println("Deleting", image, ":", tag)
		errDel := registry.DeleteImage(image, tag, digest)
		if errDel != nil {
			fmt.Println(errDel.Error())
		}
		results <- tag
	}
}

func exit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
