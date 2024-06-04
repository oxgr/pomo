package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Model struct {
	Start    int  `json:"start"`
	Duration int  `json:"duration"`
	Active   bool `json:"active"`
	Done     bool `json:"done"`
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

		if model.Done {
			fmt.Println("timer is done!")
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

func info() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	modelJson, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(modelJson))

	return nil
}
