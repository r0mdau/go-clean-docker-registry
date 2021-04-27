package main

import (
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"testing"
)

func TestInitCmdAppConfiguration(t *testing.T) {
	t.Run("Configuration values", func(t *testing.T) {
		app := newTestApp()
		expect(t, app.Name, "Go Clean Docker Registry")
		expect(t, app.EnableBashCompletion, true)
	})
}

func TestCommandShowimagesAppValues(t *testing.T) {
	var registyUrl string
	var insecure bool

	t.Run("Showimages test flags default secure", func(t *testing.T) {
		app := newTestApp()
		// ugly, don't change orders in cmd.go
		app.Commands[0].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			insecure = c.Bool("insecure")
			return nil
		}

		err := app.Run([]string{"", "showimages", "--url", "https://example.com"})

		expect(t, registyUrl, "https://example.com")
		expect(t, insecure, false)
		assertNoError(t, err)
	})

	t.Run("Showimages test flags insecure", func(t *testing.T) {
		app := newTestApp()
		app.Commands[0].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			insecure = c.Bool("insecure")
			return nil
		}

		err := app.Run([]string{"", "showimages", "--url", "https://example.com", "--insecure"})

		expect(t, registyUrl, "https://example.com")
		expect(t, insecure, true)
		assertNoError(t, err)
	})
}

func TestCommandShowtagsAppValues(t *testing.T) {
	var registyUrl, image string
	var insecure bool

	t.Run("Showtags test minimum flags", func(t *testing.T) {
		app := newTestApp()
		app.Commands[1].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			image = c.String("image")
			return nil
		}

		err := app.Run([]string{"", "showtags", "--url", "https://example.com", "--image", "test"})

		expect(t, registyUrl, "https://example.com")
		expect(t, image, "test")
		assertNoError(t, err)
	})
	t.Run("Showtags test maximum flags", func(t *testing.T) {
		app := newTestApp()
		app.Commands[1].Action = func(c *cli.Context) error {
			registyUrl = c.String("url")
			image = c.String("image")
			insecure = c.Bool("insecure")
			return nil
		}

		err := app.Run([]string{"", "showtags", "--url", "https://example.com", "--image", "test", "--insecure"})

		expect(t, registyUrl, "https://example.com")
		expect(t, image, "test")
		expect(t, insecure, true)
		assertNoError(t, err)
	})
}

func TestCommandShowimagesRequiredFlagAppRunBehavior(t *testing.T) {
	tdata := []struct {
		testCase        string
		appRunInput     []string
		expectedAnError bool
	}{
		{
			testCase:        "valid_case_empty_input",
			appRunInput:     []string{"myCLI"},
			expectedAnError: false,
		},
		{
			testCase:        "error_case_empty_input_with_required_flag_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages"},
			expectedAnError: true,
		},
		{
			testCase:        "valid_case_with_minimum_required_flag_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost"},
			expectedAnError: false,
		},
		{
			testCase:        "valid_case_with_maximum_required_flag_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost", "--insecure"},
			expectedAnError: false,
		},
		{
			testCase:        "error_case_not_allowed_tag_image_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost", "--image", "r0mdau/nodejs"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_not_allowed_tag_flag_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost", "--tag", "1.0.0"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_not_allowed_keep_flag_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost", "--keep", "1.0.0"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_not_allowed_keep_dryrun_on_command_showimages",
			appRunInput:     []string{"myCLI", "showimages", "--url", "http://localhost", "--dryrun"},
			expectedAnError: true,
		},
	}

	assertAppBehaviour(t, tdata)
}

func TestCommandShowtagsRequiredFlagAppRunBehavior(t *testing.T) {
	tdata := []struct {
		testCase        string
		appRunInput     []string
		expectedAnError bool
	}{
		{
			testCase:        "valid_case_empty_input",
			appRunInput:     []string{"myCLI"},
			expectedAnError: false,
		},
		{
			testCase:        "error_case_empty_input_with_required_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_missing_url_required_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--image", "r0mdau/nodejs"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_missing_image_required_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost"},
			expectedAnError: true,
		},
		{
			testCase:        "valid_case_with_minimum_required_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost", "--image", "r0mdau/nodejs"},
			expectedAnError: false,
		},
		{
			testCase:        "valid_case_with_maximum_required_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost", "--image", "r0mdau/nodejs", "--insecure"},
			expectedAnError: false,
		},
		{
			testCase:        "error_case_not_allowed_tag_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost", "--image", "r0mdau/nodejs", "--tag", "1.0.0"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_not_allowed_keep_flag_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost", "--image", "r0mdau/nodejs", "--keep", "1.0.0"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_not_allowed_keep_dryrun_on_command_showtags",
			appRunInput:     []string{"myCLI", "showtags", "--url", "http://localhost", "--image", "r0mdau/nodejs", "--dryrun"},
			expectedAnError: true,
		},
	}

	assertAppBehaviour(t, tdata)
}

func TestCommandDeleteRequiredFlagAppRunBehavior(t *testing.T) {
	tdata := []struct {
		testCase        string
		appRunInput     []string
		expectedAnError bool
	}{
		{
			testCase:        "valid_case_empty_input",
			appRunInput:     []string{"myCLI"},
			expectedAnError: false,
		},
		{
			testCase:        "error_case_empty_input_with_required_flag_on_command_delete",
			appRunInput:     []string{"myCLI", "delete"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_missing_url_required_flag_on_command_delete",
			appRunInput:     []string{"myCLI", "delete", "--image", "r0mdau/nodejs"},
			expectedAnError: true,
		},
		{
			testCase:        "error_case_missing_image_required_flag_on_command_delete",
			appRunInput:     []string{"myCLI", "delete", "--url", "http://localhost"},
			expectedAnError: true,
		},
		{
			testCase:        "valid_case_with_minimum_required_flag_on_command_delete",
			appRunInput:     []string{"myCLI", "delete", "--url", "http://localhost", "--image", "r0mdau/nodejs"},
			expectedAnError: false,
		},
		{
			testCase:        "valid_case_with_maximum_required_flag_on_command_delete",
			appRunInput:     []string{"myCLI", "delete", "--url", "http://localhost", "--image", "r0mdau/nodejs", "--tag", "1.0.0", "--keep", "1", "--dryrun", "--insecure"},
			expectedAnError: false,
		},
	}

	assertAppBehaviour(t, tdata)
}

func assertAppBehaviour(t *testing.T, tdata []struct {
	testCase        string
	appRunInput     []string
	expectedAnError bool
}) {
	for _, test := range tdata {
		t.Run(test.testCase, func(t *testing.T) {
			t.Helper()
			app := newTestApp()
			err := app.Run(test.appRunInput)

			if test.expectedAnError && err == nil {
				t.Errorf("expected an error, but there was none")
			}
			if !test.expectedAnError && err != nil {
				t.Errorf("did not expected an error, but there was one: %s", err)
			}
		})
	}
}

func newTestApp() *cli.App {
	app := initCmdApp()
	for _, command := range app.Commands {
		command.Action = func(c *cli.Context) error {
			return nil
		}
	}
	app.Writer = ioutil.Discard
	return app
}
