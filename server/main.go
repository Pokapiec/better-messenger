package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Server struct {
	connections map[*websocket.Conn]bool
}

func (s *Server) AddConn(newConn *websocket.Conn) {
	s.connections[newConn] = true
}

func (s *Server) RemoveConn(conn *websocket.Conn) {
	delete(s.connections, conn)
}

func (s *Server) WSSHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
		}
		log.Println("Upgraded connection.")
		s.AddConn(conn)
		log.Println("Added connection.")

		defer conn.Close()

		for {
			var jsonMsg WSMessage
			err := conn.ReadJSON(&jsonMsg)
			if err != nil {
				log.Println("Error reading msg:", err)
				s.RemoveConn(conn)
				return
			}
			fmt.Println("received msg:", jsonMsg)

			for c := range s.connections {
				if c == conn {
					continue
				}
				err = c.WriteJSON(jsonMsg)
				if err != nil {
					log.Println("Failed writing msg:", err)
					return
				}
			}
		}
	}

}

func HanlderMessageList(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversationIdStr := r.PathValue("id")
		conversationId, err := strconv.Atoi(conversationIdStr)
		if err != nil {
			http.Error(w, "invalid conversation id", http.StatusBadRequest)
			return
		}

		messages, err := storage.GetMessages(conversationId)
		if err != nil {
			http.Error(w, "Error fetching data", http.StatusBadRequest)
			return
		}

		log.Println("retrieved messages:", messages)

		JSONResponse(w, HttpResponse{data: messages})
	}
}

func main() {
	server := Server{connections: make(map[*websocket.Conn]bool)}
	router := http.NewServeMux()
	storage, err := NewStorage()
	if err != nil {
		log.Fatalln("Failed to initialize storage:", err)
	}
	router.HandleFunc("GET /ws", server.WSSHandler(storage))
	router.HandleFunc("GET /conversations/{id}/messages", HanlderMessageList(storage))

	log.Println("Listening on :3001...")
	http.ListenAndServe("127.0.0.1:3001", router)
}
