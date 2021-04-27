package main

import (
	"flag"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"testing"
)

func TestConfigureRegistry(t *testing.T) {
	t.Run("TODO Cannot set cli.flagSet so configure registry return empty registry", func(t *testing.T) {
		app := &cli.App{Writer: ioutil.Discard}
		set := flag.NewFlagSet("test", 0)
		testArgs := []string{"", "showimages", "--url", "https://example.com"}
		set.Parse(testArgs)
		context := cli.NewContext(app, set, nil)

		command := cli.Command{
			Name:            "showimages",
			Usage:           "this is for testing",
			Description:     "testing",
			Action:          func(_ *cli.Context) error { return nil },
			SkipFlagParsing: true,
		}

		command.Run(context)

		rActual := configureRegistry(context)

		rExpected := Registry{}
		rExpected.configure("", false)
		expect(t, rActual.BaseUrl, rExpected.BaseUrl)
		expect(t, "todo", "todo")
	})
}
