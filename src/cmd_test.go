package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func TestInitCmdApp(t *testing.T) {
	var url, image, tag string
	var keep int
	var dryrun, insecure bool

	t.Run("Configuration values", func(t *testing.T) {
		app := initCmdApp()
		equals(t, "Go Clean Docker Registry", app.Name)
		equals(t, true, app.EnableBashCompletion)
	})

	t.Run("Minimum flags", func(t *testing.T) {
		app := initCmdApp()
		app.Action = func(c *cli.Context) error {
			url = c.String("url")
			image = c.String("image")
			tag = c.String("tag")
			keep = c.Int("keep")
			dryrun = c.Bool("dryrun")
			insecure = c.Bool("insecure")
			return nil
		}

		_ = app.Run([]string{"", "--url", "https://example.com", "--image", "test"})

		equals(t, "https://example.com", url)
		equals(t, "test", image)
		equals(t, "", tag)
		equals(t, 0, keep)
		equals(t, false, dryrun)
		equals(t, false, insecure)
	})
	t.Run("Maximum flags", func(t *testing.T) {
		app := initCmdApp()
		app.Action = func(c *cli.Context) error {
			url = c.String("url")
			image = c.String("image")
			tag = c.String("tag")
			keep = c.Int("keep")
			dryrun = c.Bool("dryrun")
			insecure = c.Bool("insecure")
			return nil
		}

		_ = app.Run([]string{"", "--url", "https://example.com", "--image", "test", "--tag", "master", "--keep", "5", "--dryrun", "--insecure"})

		equals(t, "https://example.com", url)
		equals(t, "test", image)
		equals(t, "master", tag)
		equals(t, 5, keep)
		equals(t, true, dryrun)
		equals(t, true, insecure)
	})
}
