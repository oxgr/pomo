package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start a timer for <n> minutes",
				Action: func(cCtx *cli.Context) error {
					arg := cCtx.Args().First()

					mins, err := func(str string) (int, error) {
						if str == "" {
							return 0, nil
						}

						mins, err := strconv.Atoi(arg)
						if err != nil {
							fmt.Println("could not parse number!")
							return 0, err
						}

						return mins, nil
					}(arg)
					if err != nil {
						return err
					}

					return start(mins)
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"x"},
				Usage:   "stop the timer and scrap it",
				Action: func(cCtx *cli.Context) error {
					return stop()
				},
			},
			{
				Name:    "pause",
				Aliases: []string{"p"},
				Usage:   "pause the timer, keeping old values",
				Action: func(cCtx *cli.Context) error {
					return pause()
				},
			},
			{
				Name:    "show",
				Aliases: []string{"t"},
				Usage:   "show the status or remaining time",
				Action: func(cCtx *cli.Context) error {
					return show()
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "show info",
				Action: func(cCtx *cli.Context) error {
					return info()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
