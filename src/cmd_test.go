package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func TestInitCmdAppConfiguration(t *testing.T) {
	t.Run("Configuration values", func(t *testing.T) {
		app := initCmdApp()
		equals(t, "Go Clean Docker Registry", app.Name)
		equals(t, true, app.EnableBashCompletion)
	})
}

func TestInitCmdAppShowimages(t *testing.T) {
	var registyUrl string
	var insecure bool

	t.Run("Showimages test flags default secure", func(t *testing.T) {
		app := initCmdApp()
		// ugly, don't change orders in cmd.go
		app.Commands[0].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			insecure = c.Bool("insecure")
			return nil
		}

		_ = app.Run([]string{"", "showimages", "--url", "https://example.com"})

		equals(t, "https://example.com", registyUrl)
		equals(t, false, insecure)
	})

	t.Run("Showimages test flags insecure", func(t *testing.T) {
		app := initCmdApp()
		app.Commands[0].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			insecure = c.Bool("insecure")
			return nil
		}

		_ = app.Run([]string{"", "showimages", "--url", "https://example.com", "--insecure"})

		equals(t, "https://example.com", registyUrl)
		equals(t, true, insecure)
	})
}

func TestInitCmdAppShowtags(t *testing.T) {
	var registyUrl, image string
	var insecure bool

	t.Run("Showtags test minimum flags", func(t *testing.T) {
		app := initCmdApp()
		app.Commands[1].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			image = c.String("image")
			return nil
		}

		_ = app.Run([]string{"", "showtags", "--url", "https://example.com", "--image", "test"})

		equals(t, "https://example.com", registyUrl)
		equals(t, "test", image)
	})
	t.Run("Shotags test maximum flags", func(t *testing.T) {
		app := initCmdApp()
		app.Commands[1].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			image = c.String("image")
			insecure = c.Bool("insecure")
			return nil
		}

		_ = app.Run([]string{"", "showtags", "--url", "https://example.com", "--image", "test", "--insecure"})

		equals(t, "https://example.com", registyUrl)
		equals(t, "test", image)
		equals(t, true, insecure)
	})
}
