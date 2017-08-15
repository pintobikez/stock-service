package main

import (
	"fmt"
	"github.com/labstack/gommon/color"
	//errors "github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

// Check the health of config file and ssl certificates
func Check(c *cli.Context) error {

	errs := []error{}
	success := []string{}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Printf("%s %s\n", color.Red("[ERROR]"), err.Error())
		}

		cli.OsExiter(1)
	}

	for _, msg := range success {
		fmt.Printf("%s %s\n", color.Green("[LOADED]"), msg)
	}

	return nil
}
