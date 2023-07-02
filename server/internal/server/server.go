package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kyosu-1/reversi/server/internal/game"
	"github.com/kyosu-1/reversi/server/internal/strategy"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Server struct {
	Game      *game.Board
	Clients   map[*websocket.Conn]bool
	Broadcast chan Message
	Strategy  strategy.Strategy
}

func NewServer(game *game.Board, strategy strategy.Strategy) *Server {
	return &Server{
		Game:      game,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan Message),
		Strategy:  strategy,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) run() {
	for {
		select {
		case message := <-s.Broadcast:
			for client := range s.Clients {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(s.Clients, client)
				}
			}
		}
	}
}

func ServeWs(s *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.Clients[conn] = true
	go s.handleMessages(conn)
	go s.run()

	msg := Message{"init", s.Game}
	conn.WriteJSON(msg)
}

func (s *Server) handleMessages(conn *websocket.Conn) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			conn.Close()
			delete(s.Clients, conn)
			return
		}

		switch msg.Type {
		case "move":
			move := msg.Data.(map[string]interface{})
			x := int(move["x"].(float64))
			y := int(move["y"].(float64))
			err := s.Game.MakeMove(game.Move{X: x, Y: y})
			if err == nil && s.Game.Turn == game.White {
				move := s.Strategy.BestMove(s.Game)
				s.Game.MakeMove(move)
			}
		}

		s.Broadcast <- Message{"update", s.Game}
	}
}
