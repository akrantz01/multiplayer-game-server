package main

import (
	"encoding/json"
	"sync"
)

type (
	// Server info
	Server struct {
		Mode	  int	   `yaml:"mode"`
		Host      string   `yaml:"host"`
		Port      string   `yaml:"port"`
		Debug     bool	   `yaml:"debug"`
		StaticDir string   `yaml:"static-directory"`
	}

	// Game data
	GameData struct {
		sync.RWMutex
		Globals		map[string]map[string]Value
		Users		map[string]*UserValue
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

	// Test player info
	TestPlayer struct {
		Starting	Coordinates		`yaml:"starting"`
		Movements	movement		`yaml:"movement"`
		Reference	*UserValue
	}

	// Movement for test players
	movement struct {
		Active		bool		`yaml:"active"`
		Actions		[]action	`yaml:"actions"`
	}

	// Actions for test players' movement
	action struct {
		From	float32		`yaml:"from"`
		To		float32		`yaml:"to"`
		Axis	int		`yaml:"axis"`
		Speed	float32	`yaml:"speed"`
	}
)

func (g *GameData) SetGlobals(globals map[string]map[string]Value) {
	g.Lock()
	defer g.Unlock()
	g.Globals = globals
}

func (g *GameData) SetUserData(id string, x, y, z, orientation float32, other map[string]interface{}) {
	g.Lock()
	defer g.Unlock()
	g.Users[id] = &UserValue{
		X: x,
		Y: y,
		Z: z,
		Orientation: orientation,
		Other: other,
	}
}

func (g *GameData) SetObject(id string, x, y, z float32, other map[string]interface{}) {
	g.Lock()
	defer g.Unlock()
	g.Objects[id] = Object{
		Coordinates: Coordinates{
			X: x,
			Y: y,
			Z: z,
		},
		Other: other,
	}
}

func (g *GameData) CreateTestUser(id string, x, y, z float32) *UserValue {
	g.Lock()
	defer g.Unlock()
	g.Users[id] = &UserValue{
		X: x,
		Y: y,
		Z: z,
		Orientation: 0,
		Other: make(map[string]interface{}),
	}
	return g.Users[id]
}

func (g *GameData) DeleteUserData(id string) {
	g.Lock()
	defer g.Unlock()
	delete(g.Users, id)
}

func (g *GameData) DeleteObject(id string) {
	g.Lock()
	defer g.Unlock()
	delete(g.Objects, id)
}

func (g *GameData) GetAllData() ([]byte, error) {
	g.Lock()
	defer g.Unlock()
	j, err := json.Marshal(struct {
		Globals		map[string]map[string]Value
		Users		map[string]*UserValue
		Objects		map[string]Object
	}{
		Globals: g.Globals,
		Users: g.Users,
		Objects: g.Objects,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (u UserValue) equals (u2 UserValue) bool {
	if u.X != u2.X {return false}
	if u.Y != u2.Y {return false}
	if u.Z != u2.Z {return false}
	if len(u.Other) != len(u2.Other) {return false}
	return true
}

func (t *TestPlayer) Move (id int) {
	for i, a := range t.Movements.Actions {
		switch a.Axis {
		case 0:
			t.Reference.X += a.Speed
			break
		case 1:
			t.Reference.Y += a.Speed
			break
		case 2:
			t.Reference.Z += a.Speed
			break
		}

		if (t.Reference.X <= a.From || t.Reference.X >= a.To) && a.Axis == 0 {
			t.Movements.Actions[i].Speed = -a.Speed
		} else if (t.Reference.Y <= a.From || t.Reference.Y >= a.To) && a.Axis == 1 {
			t.Movements.Actions[i].Speed = -a.Speed
		} else if (t.Reference.Z <= a.From || t.Reference.Z >= a.To) && a.Axis == 2 {
			t.Movements.Actions[i].Speed = -a.Speed
		}
	}
}
