package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSMessage struct {
	Username       string `json:"username"`
	Message        string `json:"message"`
	ConversationId int    `json:"conversation_id"`
}

type ConnectionData struct {
	Conversations []int
}

type Server struct {
	connections map[*websocket.Conn]ConnectionData
}

func (s *Server) AddConn(newConn *websocket.Conn) {
	s.connections[newConn] = ConnectionData{Conversations: []int{}}
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

			if jsonMsg.Message == "<INITIAL>" {
				value := s.connections[conn]
				value.Conversations = append(value.Conversations, jsonMsg.ConversationId)
				s.connections[conn] = value
				continue
			}

			err = storage.CreateMessage(jsonMsg)
			if err != nil {
				fmt.Println("Failed to create message:", err)
			}

			for c, cData := range s.connections {
				if c == conn {
					continue
				}

				if !slices.Contains(cData.Conversations, jsonMsg.ConversationId) {
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
			fmt.Println("error fetching messages", err)
			http.Error(w, "Error fetching data", http.StatusBadRequest)
			return
		}
		JSONResponse(w, HttpResponse{Data: messages})
	}
}

func HanlderConversationList(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversations, err := storage.GetConversations()
		if err != nil {
			http.Error(w, "Error fetching data", http.StatusBadRequest)
			return
		}
		JSONResponse(w, HttpResponse{Data: conversations})
	}
}

func main() {
	server := Server{connections: make(map[*websocket.Conn]ConnectionData)}
	router := http.NewServeMux()
	storage, err := NewStorage()
	if err != nil {
		log.Fatalln("Failed to initialize storage:", err)
	}
	router.HandleFunc("GET /ws", server.WSSHandler(storage))
	router.HandleFunc("GET /conversations/{id}/messages", HanlderMessageList(storage))
	router.HandleFunc("GET /conversations", HanlderConversationList(storage))

	log.Println("Listening on :3001...")

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}

	http.ListenAndServe("127.0.0.1:3001", corsHandler(router))
}
