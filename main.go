package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	server Server
	data = GameData{
		Users: make(map[string]*UserValue),
		Objects: make(map[string]Object),
		Globals: make(map[string]map[string]Value),
	}
	tests []TestPlayer
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections map[string]*websocket.Conn
)

func main() {
	// Get config file
	file := flag.String("config", "config.yml", "Alternative yaml configuration file")
	flag.Parse()

	// Temporary variable
	var globals map[string]map[string]Value

	// Parse config file
	server, globals, tests = ParseConfig(*file)
	data.SetGlobals(globals)

	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if server.Upstream.InUse {
		for _, loc := range server.Upstream.Locations {
			u, _ := url.Parse(loc.URL)
			targets := []*middleware.ProxyTarget{{URL: u,},}

			g := e.Group(loc.Endpoint)
			g.Use(middleware.Proxy(middleware.NewRandomBalancer(targets)))
		}
	}
	if !server.Upstream.OverrideRoot || !server.Upstream.InUse {
		e.Static("/", "./public")
	}
	if server.Mode == 1 {
		stop := make(chan struct{})
		go func() {
			for {
				for i, tp := range tests {
					tp.Move(i)
				}
				select {
				case <-time.After(500*time.Millisecond):
					break
				case <-stop:
					return
				}
			}
		}()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			stop<-struct{}{}
			os.Exit(0)
		}()
	}

	// Define all routes
	e.GET("/ws", wsHandler)

	// Debug routes
	debug := e.Group("/debug")
	debug.Use(middleware.BasicAuth(func(u, p string, ctx echo.Context) (bool, error) {
		if u == server.DebugUser && p == server.DebugPass {
			return true, nil
		}
		return false, nil
	}))
	debug.GET("/", debugHandler)

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
				data.DeleteUserData(userID)
				delete(connections, userID)
				break
			} else {
				panic(err)
			}
		}

		connections[userID] = ws

		switch msg.Type {
		case 1:
			userID = msg.ID

			data.SetUserData(
				msg.ID,
				msg.Coordinates.X,
				msg.Coordinates.Y,
				msg.Coordinates.Z,
				msg.Orientation,
				msg.Other,
			)
			break

		case 2:
			data.SetObject(
				msg.ID,
				msg.Coordinates.X,
				msg.Coordinates.Y,
				msg.Coordinates.Z,
				msg.Other,
			)
			break

		case 3:
			err = ws.WriteJSON(data.GetAllData())
			break

		case 4:
			data.DeleteObject(msg.ID)
			break

		case 5:
			go broadcast(msg)
		}

		if err != nil {
			if strings.Contains(err.Error(), "close 1001") {
				data.DeleteUserData(userID)
				delete(connections, userID)
				break
			} else {
				panic(err)
			}
		}
	}

	return nil
}

func debugHandler(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, data.GetAllData())
}

func broadcast(msg Message) {
	for _, ws := range connections {
		if err := ws.WriteJSON(&msg); err != nil {
			panic(err)
		}
	}
}
