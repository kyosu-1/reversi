package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kyosu-1/reversi/server/internal/game"
	"github.com/kyosu-1/reversi/server/internal/server"
	"github.com/kyosu-1/reversi/server/internal/strategy"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	s := server.NewServer(game.NewBoard(), &strategy.RandomStrategy{})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(s, w, r)
	})

	log.Printf("Starting server on %s", *addr)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
