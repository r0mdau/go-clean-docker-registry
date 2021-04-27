package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func configureRegistry(c *cli.Context) Registry {
	registry := Registry{}
	registry.configure(c.String("url"), c.Bool("insecure"))
	return registry
}

func printRegistryCatalog(c *cli.Context) error {
	registry := configureRegistry(c)
	catalog := registry.getCatalog("/v2/_catalog")

	fmt.Println(string(catalog))
	return nil
}

func printRegistryTags(c *cli.Context) error {
	registry := configureRegistry(c)
	registryResponse := registry.getTagsList("/v2/" + c.String("image") + "/tags/list")

	fmt.Println(string(registryResponse.Body))
	fmt.Println("Total of", len(registryResponse.getRegistryImage().Tags), "tags.")
	return nil
}

func deleteRegistryTags(c *cli.Context) error {
	registry := configureRegistry(c)
	registryResponse := registry.getTagsList("/v2/" + c.String("image") + "/tags/list")

	fmt.Println("Total of", len(registryResponse.getRegistryImage().Tags), "tags.")
	for _, value := range registryResponse.getRegistryImage().Tags {
		if c.Bool("dryrun") {
			fmt.Println("Deleting :", value)
		}
	}
	return nil
}
