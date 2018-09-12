package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ParseConfig(file string) (Server, map[string]map[string]Value) {
	// Open file
	opened, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer opened.Close()

	// Parse json into struct
	var config struct{
		Server	Server                       `json:"server"`
		Globals	map[string]map[string]Value `json:"globals"`
	}
	content, _ := ioutil.ReadAll(opened)
	yaml.Unmarshal(content, &config)

	return config.Server, config.Globals
}
