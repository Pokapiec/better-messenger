package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var dbSchema = `
CREATE TABLE IF NOT EXISTS conversations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT
);

CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message TEXT,
	user_id INTEGER,
	conversation_id INTEGER
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT
);
`

type Message struct {
	Id             int    `db:"id" json:"id"`
	Message        string `db:"message" json:"message"`
	UserId         int    `db:"user_id" json:"user_id"`
	ConversationId int    `db:"conversation_id" json:"conversation_id"`
}

type Conversation struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
}

type Storage struct {
	DB *sqlx.DB
}

type MessageForFrontend struct {
	Id             int    `db:"id" json:"id"`
	Message        string `db:"message" json:"message"`
	UserId         int    `db:"user_id" json:"user_id"`
	Username       string `db:"username" json:"username"`
	ConversationId int    `db:"conversation_id" json:"conversation_id"`
}

func NewStorage() (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	db.MustExec(dbSchema)
	return &Storage{DB: db}, nil
}

func (s Storage) GetMessages(conversationId int) ([]MessageForFrontend, error) {
	var messages []MessageForFrontend
	err := s.DB.Select(&messages, `SELECT 
	messages.id, message, user_id, users.username, conversation_id 
	FROM messages 
	LEFT JOIN users ON messages.user_id = users.id
	WHERE conversation_id = $1`,
		conversationId,
	)
	if err != nil {
		return []MessageForFrontend{}, err
	}

	return messages, nil
}

func (s Storage) GetConversations() ([]Conversation, error) {
	var conversations []Conversation
	err := s.DB.Select(&conversations, "SELECT id, name FROM conversations")
	if err != nil {
		log.Println("Error when fetching conversations:", err)
		return []Conversation{}, err
	}

	return conversations, nil
}

func (s Storage) GetOrCreateUser(username string) (User, error) {
	user := User{}

	err := s.DB.Get(&user, "SELECT id, username FROM users WHERE username = $1", username)
	if err == nil {
		fmt.Println("User for username already exists.")
		return user, nil
	}
	fmt.Println("Creating user...")

	var insertedId int
	err = s.DB.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING id", username).Scan(&insertedId)
	user = User{Id: insertedId, Username: username}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s Storage) CreateMessage(msgJson WSMessage) error {
	user, err := s.GetOrCreateUser(msgJson.Username)
	if err != nil {
		return err
	}

	_, err = s.DB.NamedExec("INSERT INTO messages (message, user_id, conversation_id) VALUES (:message, :user_id, :conversation_id)", map[string]interface{}{
		"message":         msgJson.Message,
		"user_id":         user.Id,
		"conversation_id": msgJson.ConversationId,
	})
	if err != nil {
		return err
	}
	return nil
}
