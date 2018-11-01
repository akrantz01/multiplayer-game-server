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
		data.Users[strconv.Itoa(id)] = &UserValue{
			X: player.Starting.X,
			Y: player.Starting.Y,
			Z: player.Starting.Z,
			Orientation: 0,
			Other: make(map[string]interface{}),
		}
		p := data.Users[strconv.Itoa(id)]
		config.TestPlayers[id].Reference = p
	}

	return config.Server, config.Globals, config.TestPlayers
}
