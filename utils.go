package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
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
