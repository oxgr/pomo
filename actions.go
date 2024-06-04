package main

import (
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

func show(formatJson bool) error {
	model, err := readModel()
	if err != nil {
		return err
	}

	type Show struct {
		Text  string `json:"text"`
		Class string `json:"class"`
	}

	if !model.Active {
		showObj := Show{}

		if model.Duration == 0 {
			showObj.Class = "stopped"
		} else if !model.Done {
			showObj.Class = "paused"
		} else {
			showObj.Class = "done"
		}

		if !formatJson {
			fmt.Println(showObj.Class)
			return nil
		}

		jsonStr, err := toJson(showObj, false)
		if err != nil {
			return err
		}

		fmt.Println(jsonStr)
		return nil
	}

	remaining := model.Start + model.Duration - int(time.Now().Unix())

	if remaining < 0 {
		model.Active = false
		model.Done = true
		if err := writeModel(model); err != nil {
			return err
		}

		showObj := Show{
			Class: "done",
		}

		if !formatJson {
			fmt.Println(showObj.Class)
			return nil
		}

		jsonStr, err := toJson(showObj, false)
		if err != nil {
			return err
		}

		fmt.Println(jsonStr)
		return nil
	}

	rem_mins := remaining / 60
	rem_secs := remaining % 60
	showObj := Show{
		Text:  fmt.Sprintf("%02d:%02d", rem_mins, rem_secs),
		Class: "active",
	}

	if !formatJson {
		fmt.Println(showObj.Text)
		return nil
	}

	jsonStr, err := toJson(showObj, false)
	if err != nil {
		return err
	}

	fmt.Println(jsonStr)
	return nil
}

func info() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	fmt.Println(toJson(model, true))

	return nil
}
