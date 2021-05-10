package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/r0mdau/go-clean-docker-registry/internal/filter"
	"github.com/r0mdau/go-clean-docker-registry/pkg/registry"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
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

	fmt.Println(string(repositories.Body))
	fmt.Fprintf(os.Stderr, "Total of %d repositories.\n", len(repositories.GetRepository().List))
	return nil
}

func printImageTagsList(c *cli.Context) error {
	registry := registry.NewRegistry(c.String("url"), c.Bool("insecure"))
	verifyRegistryVersion(registry)

	imageTags, err := registry.ListImageTags(c.String("image"))
	exit(err)

	fmt.Println(string(imageTags.Body))
	fmt.Fprintf(os.Stderr, "Total of %d tags.\n", len(imageTags.GetImage().Tags))
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
		output, _ := json.Marshal(tagsToDelete)
		fmt.Println(string(output))
		fmt.Fprintf(os.Stderr, "Dryrun, it should delete image : \"%s\" with %d tags.\n", cliImage, len(tagsToDelete))
		return nil
	}

	if confirm("Are you sure to delete these tags ? (maybe try --dryrun first)") {
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
		fmt.Fprintf(os.Stderr, "Total of %d tags deleted.\n", len(tagsToDelete))
	}
	return nil
}

func wDelete(registry registry.Registry, image string, jobs <-chan string, results chan<- string) {
	for tag := range jobs {
		digest, errGet := registry.GetDigestFromManifest(image, tag)
		if errGet != nil {
			fmt.Fprintf(os.Stderr, "%s\n", errGet.Error())
			results <- tag
			continue
		}
		fmt.Fprintf(os.Stderr, "Deleting %s:%s\n", image, tag)
		errDel := registry.DeleteImage(image, tag, digest)
		if errDel != nil {
			fmt.Fprintf(os.Stderr, "%s\n", errGet.Error())
		}
		results <- tag
	}
}

func confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
		fmt.Fprintf(os.Stderr, "Canceled.\n")
		return false
	}
}

func exit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
