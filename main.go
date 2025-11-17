package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	games map[string]*Game
	mu    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		games: make(map[string]*Game),
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	var initMsg struct {
		Type     string `json:"type"`
		GameID   string `json:"gameId"`
		PlayerID string `json:"playerId"`
		Name     string `json:"name"`
	}

	if err := conn.ReadJSON(&initMsg); err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	s.mu.Lock()
	game, exists := s.games[initMsg.GameID]
	if !exists {
		game = NewGame(initMsg.GameID)
		s.games[initMsg.GameID] = game
		go game.Run()
	}
	s.mu.Unlock()

	player := &Player{
		ID:   initMsg.PlayerID,
		Name: initMsg.Name,
		Conn: conn,
		Game: game,
	}

	game.Register <- player

	// Read messages from this player
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read error: %v", err)
			game.Unregister <- player
			break
		}
		msg.PlayerID = player.ID
		game.Broadcast <- msg
	}
}

func (s *Server) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	gameID := uuid.New().String()

	s.mu.Lock()
	s.games[gameID] = NewGame(gameID)
	go s.games[gameID].Run()
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"gameId": gameID,
	})
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
	} else {
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/ws", server.handleWebSocket)
	http.HandleFunc("/api/create", server.handleCreateGame)
	http.HandleFunc("/", server.handleStatic)

	port := ":8080"
	fmt.Printf("🚀 Space Adventure Server starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
