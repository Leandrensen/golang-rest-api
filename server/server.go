package server

import (
	"context"
	"errors"
	"golang-rest-api-websockets/database"
	"golang-rest-api-websockets/repository"
	"golang-rest-api-websockets/websocket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

// Server interface without Websockets
// type Server interface {
// 	Config() *Config
// }

// Server interface with Websockets
type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

// Broker without WebSockets
// type Broker struct {
// 	config *Config
// 	router *mux.Router
// }

// Broker with Websockets
type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

// Needed for Websocket implementation
func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

// ctx context.Context is needed to know where an error happened (Because we will be having a lot of go routines)
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("database is required")
	}

	// Broker without Websockets
	// broker := &Broker{
	// 	config: config,
	// 	router: mux.NewRouter(),
	// }

	// Broker with Websockets
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	// b.router = mux.NewRouter()
	// ^^^ It creates again the router that was created on line 47
	// I belive it's an error
	binder(b, b.router)
	// This handler is created to stop the CORS
	// when accessing from a website vvvv
	handler := cors.Default().Handler(b.router)
	repo, err := database.NewPostgresrepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	// Websockets Hub vvv
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
