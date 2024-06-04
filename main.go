package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

type Model struct {
	Start    int  `json:"start"`
	Duration int  `json:"duration"`
	Active   bool `json:"active"`
	Done     bool `json:"done"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start a timer for <n> minutes.",
				Action: func(cCtx *cli.Context) error {
					arg := cCtx.Args().First()

					mins := func(str string) int {
						if str != "" {
							mins, err := strconv.Atoi(arg)
							check(err)
							return mins
						}
						return 0
					}(arg)

					return start(mins)
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"S"},
				Usage:   "stop the timer and scrap it.",
				Action: func(cCtx *cli.Context) error {
					return stop()
				},
			},
			{
				Name:    "pause",
				Aliases: []string{"p"},
				Usage:   "pause the timer, keeping old values.",
				Action: func(cCtx *cli.Context) error {
					return pause()
				},
			},
			{
				Name:    "show",
				Aliases: []string{"t"},
				Usage:   "show the status - remaining time or break!",
				Action: func(cCtx *cli.Context) error {
					return show()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getDataPath() string {
	const FALLBACK_FILE_PATH = "/tmp/pomo.json"

	cache_path, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("cannot get cache dir. using tmp dir...")
		return FALLBACK_FILE_PATH
	}

	path := cache_path + "/pomo/data.json"
	return path
}

func writeModel(model Model) error {
	model_json, err := json.Marshal(model)
	check(err)

	model_json_string := []byte(string(model_json))
	if err := os.WriteFile(getDataPath(), model_json_string, 0o644); err != nil {
		return err
	}
	return nil
}

func readModel() (Model, error) {
	str, err := os.ReadFile(getDataPath())
	if err != nil {
		fmt.Println("no timer found!")
		return Model{}, err
	}

	var model Model

	if err := json.Unmarshal([]byte(str), &model); err != nil {
		return model, err
	}

	return model, nil
}

func start(mins int) error {
	var model Model

	if mins == 0 {
		model, err := readModel()
		if err != nil {
			return err
		}

		if model.Duration == 0 {
			fmt.Println("no time recorded - pass a <num> of mins pls!")
			return nil
		}

		if model.Active {
			fmt.Println("timer is already active!")
			return nil
		}

		model.Active = true

		if err := writeModel(model); err != nil {
			return err
		}

		fmt.Println("timer resumed...")
		return nil
	}

	seconds := mins * 60

	model = Model{
		Start:    int(time.Now().Unix()),
		Duration: seconds,
		Active:   true,
	}

	if err := writeModel(model); err != nil {
		return err
	}

	fmt.Println("timer started for", mins, "minutes...")
	return nil
}

func stop() error {
	model := Model{}

	if err := writeModel(model); err != nil {
		return err
	}

	fmt.Println("timer stopped!")
	return nil
}

func pause() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	model.Active = false

	if err := writeModel(model); err != nil {
		return err
	}

	fmt.Println("timer paused~")

	return nil
}

func show() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	if model.Duration == 0 {
		fmt.Println("stopped")
		return nil
	}

	if !model.Active {
		fmt.Println("paused")
		return nil
	}

	remaining := model.Start + model.Duration - int(time.Now().Unix())

	if remaining < 0 {
		fmt.Println("done")
		model.Active = false
		model.Done = true
		if err := writeModel(model); err != nil {
			return err
		}
		return nil
	}

	rem_mins := remaining / 60
	rem_secs := remaining % 60
	fmt.Printf("%02d:%02d\n", rem_mins, rem_secs)
	return nil
}
