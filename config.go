package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

func ParseConfig(file string) (Server, map[string]map[string]Value, []TestPlayer) {
	// Open file
	opened, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer opened.Close()

	// Parse json into struct
	var config struct{
		Server			Server                       `json:"server"`
		Globals			map[string]map[string]Value  `json:"globals"`
		TestPlayers		[]TestPlayer				 `yaml:"test-players"`
	}
	content, _ := ioutil.ReadAll(opened)
	yaml.Unmarshal(content, &config)

	for id, player := range config.TestPlayers {
		config.TestPlayers[id].Reference = data.CreateTestUser(
			strconv.Itoa(id),
			player.Starting.X,
			player.Starting.Y,
			player.Starting.Z,
		)
	}

	return config.Server, config.Globals, config.TestPlayers
}
