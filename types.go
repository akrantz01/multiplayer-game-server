package main

import "sync"

type (
	// Server info
	Server struct {
		Host      string   `yaml:"host"`
		Port      string   `yaml:"port"`
		DebugUser string   `yaml:"debug-username"`
		DebugPass string   `yaml:"debug-password"`
		Upstream  upstream `yaml:"upstream"`
	}

	// upstream info
	upstream struct {
		InUse			bool    	`yaml:"active"`
		OverrideRoot	bool		`yaml:"override-root"`
		Locations		[]location	`yaml:"locations"`
	}

	// location info
	location struct {
		URL			string	`yaml:"url"`
		Endpoint	string	`yaml:"endpoint"`
	}

	// Game data
	GameData struct {
		sync.Mutex
		Globals		map[string]map[string]Value
		Users		map[string]UserValue
		Objects		map[string]Object
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
		Z		float32
		Orientation	float32
		Other	map[string]interface{}
	}

	// Message container
	Message struct {
		Type		int						`yaml:"type"`
		ID			string					`yaml:"id"`
		Other		map[string]interface{}	`yaml:"other"`
		Coordinates Coordinates				`yaml:"coordinates"`
		Orientation	float32					`yaml:"orientation"`
	}

	// Store player coordinates (2d)
	Coordinates struct {
		X	float32	`yaml:"x"`
		Y	float32 `yaml:"y"`
		Z	float32	`yaml:"z"`
	}

	// Store objects
	Object struct {
		Coordinates Coordinates				`yaml:"coordinates"`
		Other		map[string]interface{}	`yaml:"other"`
	}
)

func (u UserValue) equals (u2 UserValue) bool {
	if u.X != u2.X {return false}
	if u.Y != u2.Y {return false}
	if u.Z != u2.Z {return false}
	if len(u.Other) != len(u2.Other) {return false}
	return true
}
