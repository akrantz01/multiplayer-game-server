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
	"sync"
	"syscall"
	"time"
)

var (
	server Server
	data GameData
	tests []TestPlayer
	dataMutex = sync.RWMutex{}
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

	dataMutex.Lock()
	data.Users = make(map[string]*UserValue)
	data.Objects = make(map[string]Object)
	dataMutex.Unlock()

	// Parse config file
	server, globals, tests = ParseConfig(*file)

	dataMutex.Lock()
	data.Globals = globals
	dataMutex.Unlock()

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
				delete(data.Users, userID)
				delete(connections, userID)
				break
			} else {
				panic(err)
			}
		}

		connections[userID] = ws

		//fmt.Printf("%+v", msg)

		switch msg.Type {
		case 1:
			userID = msg.ID

			dataMutex.Lock()
			data.Users[msg.ID] = &UserValue{
				X: msg.Coordinates.X,
				Y: msg.Coordinates.Y,
				Z: msg.Coordinates.Z,
				Orientation: msg.Orientation,
				Other: msg.Other,
			}
			dataMutex.Unlock()

			//fmt.Printf("%+v\n", data.Users)

			break

		case 2:
			dataMutex.Lock()
			data.Objects[msg.ID] = Object{
				Coordinates: msg.Coordinates,
				Other: msg.Other,
			}
			break

		case 3:
			dataMutex.RLock()
			err = ws.WriteJSON(&data)
			dataMutex.RUnlock()
			break

		case 4:
			dataMutex.Lock()
			delete(data.Objects, msg.ID)
			dataMutex.Unlock()
			break

		case 5:
			go broadcast(msg)
		}

		if err != nil {
			if strings.Contains(err.Error(), "close 1001") {
				delete(data.Users, userID)
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
	dataMutex.RLock()
	d := &data
	dataMutex.RUnlock()

	return ctx.JSON(http.StatusOK, d)
}

func broadcast(msg Message) {
	for _, ws := range connections {
		if err := ws.WriteJSON(&msg); err != nil {
			panic(err)
		}
	}
}
