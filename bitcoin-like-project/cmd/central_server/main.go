package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/yourusername/bitcoin-like-project/internal/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func main() {
	centralServer := server.NewCentralServer()
	go centralServer.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		centralServer.HandleConnection(conn)
	})

	mux.HandleFunc("/wallet", centralServer.HandleCreateWallet)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, 
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true, 
	})


	handler := c.Handler(mux)

	log.Println("Central server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}