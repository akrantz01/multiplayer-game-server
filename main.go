package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strings"
	"sync"
)

var (
	server Server
	data GameData
	dataMutex = sync.RWMutex{}
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	// Get config file
	file := flag.String("config", "config.yml", "Alternative yaml configuration file")
	flag.Parse()

	// Temporary variable
	var globals map[string]map[string]Value

	// Parse config file
	server, globals = ParseConfig(*file)

	dataMutex.Lock()
	data.Globals = globals
	data.Users = make(map[string]UserValue)
	dataMutex.Unlock()

	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define all routes
	e.Static("/", "./public")
	e.GET("/ws", wsHandler)

	// Debug routes
	debug := e.Group("/debug")
	debug.Use(middleware.BasicAuth(func(u, p string, ctx echo.Context) (bool, error) {
		if u == server.DebugUser && p == server.DebugPass {
			return true, nil
		}
		return false, nil
	}))
	debug.GET("/all", debugAllHandler)
	debug.GET("/globals", debugGlobalHander)
	debug.GET("/globals/", debugGlobalHander)
	debug.GET("/globals/:key", debugGlobalHander)
	debug.GET("/users", debugUserHandler)
	debug.GET("/users/", debugUserHandler)
	debug.GET("/users/:id", debugUserHandler)

	// Start server
	e.Logger.Fatal(e.Start(server.Host + ":" + server.Port))
}

func wsHandler(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	// Store user id for deletion of data
	var userID string

	// Close connection after end
	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if strings.Contains(err.Error(), "close 1001") {
				delete(data.Users, userID)
				break
			} else {
				panic(err)
			}
		}

		//fmt.Printf("%+v", msg)

		switch msg.Type {
		case 1:
			userID = msg.ID

			dataMutex.Lock()
			data.Users[msg.ID] = UserValue{
				X: msg.Coordinates.X,
				Y: msg.Coordinates.Y,
				Other: msg.Other,
			}
			dataMutex.Unlock()

			//fmt.Printf("%+v\n", data.Users)

			break

		case 2:
			dataMutex.RLock()
			err = ws.WriteJSON(&data)
			dataMutex.RUnlock()
			break
		}

		if err != nil {
			if strings.Contains(err.Error(), "close 1001") {
				delete(data.Users, userID)
				break
			} else {
				panic(err)
			}
		}
	}

	return nil
}

func debugAllHandler(ctx echo.Context) error {
	dataMutex.RLock()
	d := &data
	dataMutex.RUnlock()

	return ctx.JSON(http.StatusOK, d)
}

func debugGlobalHander(ctx echo.Context) error {
	key := ctx.Param("key")

	var g interface{}
	dataMutex.RLock()
	if data.Globals[key] == nil || key == "" {
		g = &data.Globals
	} else {
		g = data.Globals[key]
	}
	dataMutex.RUnlock()

	return ctx.JSON(http.StatusOK, g)
}

func debugUserHandler(ctx echo.Context) error {
	id := ctx.Param("id")

	var u interface{}
	dataMutex.RLock()
	if data.Users[id].equals(UserValue{}) || id == "" {
		u = &data.Users
	} else {
		u = data.Users[id]
	}
	dataMutex.RUnlock()

	return ctx.JSON(http.StatusOK, u)
}
