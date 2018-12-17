package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
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

	// Serve from static directory
	http.Handle("/", handlers.LoggingHandler(os.Stdout, http.FileServer(http.Dir(server.StaticDir))))

	// Gracefully stop goroutines
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		log.Println("Listening for SIGINT, SIGTERM & SIGKILL...")
		<-c
		log.Println("Stopped goroutines")
		os.Exit(0)
	}()

	// Enable test users
	if server.Mode == 1 {
		go func() {
			log.Println("Started moving test users...")
			for {
				for i, tp := range tests {
					tp.Move(i)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	go func() {
		log.Println("Started pushing data...")
		for {
			if out, err := data.GetAllData(); err != nil {
				fmt.Printf("JSON Encoding Error: %v", err)
			} else {
				hub.broadcast <- out
			}
			time.Sleep(15 * time.Millisecond)
		}
	}()

	// WebSocket route
	http.HandleFunc("/ws", wsHandler)

	// Debug routes
	if server.Debug { http.Handle("/debug", handlers.LoggingHandler(os.Stdout, debugHandler{})) }

	// Start server
	go hub.run()
	log.Fatal(http.ListenAndServe(server.Host + ":" + server.Port, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}

type debugHandler struct {}
func (_ debugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if j, err := data.GetAllData(); err != nil {
		log.Printf("JSON Encoding Error: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}
