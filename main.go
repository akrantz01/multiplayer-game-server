package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
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
	hub = newHub()
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
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	client := &Client{
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
	go client.pushData()

	return nil
}

func debugHandler(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, data.GetAllData())
}
