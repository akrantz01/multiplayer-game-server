package main

import "sync"

type (
	// Server info
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		DebugUser	string	`yaml:"debug-username"`
		DebugPass	string	`yaml:"debug-password"`
	}

	// Game data
	GameData struct {
		sync.Mutex
		Globals		map[string]map[string]Value
		Users		map[string]UserValue
	}

	// Generic value
	Value struct {
		Value	string	`yaml:"value"`
		Type	string	`yaml:"type"`
	}

	// User value
	UserValue struct {
		X		float32
		Y	 	float32
		Other	map[string]interface{}
	}

	// Message container
	Message struct {
		Type		int						`yaml:"type"`
		ID			string					`yaml:"id"`
		Other		map[string]interface{}	`yaml:"other"`
		Coordinates Coordinates				`yaml:"coordinates"`
	}

	// Store player coordinates (2d)
	Coordinates struct {
		X	float32	`yaml:"x"`
		Y	float32 `yaml:"y"`
	}
)

func (u UserValue) equals (u2 UserValue) bool {
	if u.X != u2.X {return false}
	if u.Y != u2.Y {return false}
	if len(u.Other) != len(u2.Other) {return false}
	return true
}
