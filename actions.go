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
	const DEFAULT_MINS = 25

	var model Model

	if mins == 0 {

		model, err := readModel()
		if err != nil {
			return err
		}

		if model.Active {
			fmt.Println("timer is already active!")
			return nil
		}

		mins = DEFAULT_MINS

		// if model.Duration == 0 {
		// 	fmt.Println("no time recorded - pass a <num> of mins pls!")
		// 	return nil
		// }
		//
		// if model.Done {
		// 	fmt.Println("timer is done!")
		// 	return nil
		// }
		//
		// model.Active = true
		//
		// if err := writeModel(model); err != nil {
		// 	return err
		// }
		//
		// fmt.Println("timer resumed...")
		// return nil
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

func done() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	model.Active = false
	model.Done = true

	if err := writeModel(model); err != nil {
		return err
	}

	fmt.Println("timer set to done~")
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
		var showObj Show

		if model.Duration == 0 {
			// showObj.Class = "stopped"
			return nil
		} else if !model.Done {
			showObj.Class = "paused"
		} else {
			showObj.Text = "done"
			showObj.Class = "done"
		}

		var str string

		if !formatJson {
			str = showObj.Class
		} else {
			str, err = toJson(showObj, false)
			if err != nil {
				return err
			}
		}

		fmt.Println(str)
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

		var str string

		if !formatJson {
			str = showObj.Class
		} else {
			str, err = toJson(showObj, false)
			if err != nil {
				return err
			}
		}

		fmt.Println(str)
		return nil
	}

	rem_mins := remaining / 60
	rem_secs := remaining % 60
	showObj := Show{
		Text:  fmt.Sprintf("%02d:%02d", rem_mins, rem_secs),
		Class: "active",
	}

	var str string

	if !formatJson {
		str = showObj.Text
	} else {
		str, err = toJson(showObj, false)
		if err != nil {
			return err
		}
	}

	fmt.Println(str)
	return nil
}

func info() error {
	model, err := readModel()
	if err != nil {
		return err
	}

	modelJson, err := toJson(model, true)
	if err != nil {
		return err
	}

	fmt.Println(modelJson)

	return nil
}
