package main

import (
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
	Id             int    `db:"id"`
	Message        string `db:"message"`
	UserId         int    `db:"user_id"`
	ConversationId int    `db:"conversation_id"`
}

type Conversation struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}

type Storage struct {
	DB *sqlx.DB
}

func NewStorage() (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	db.MustExec(dbSchema)
	return &Storage{DB: db}, nil
}

func (s Storage) GetMessages(conversationId int) ([]Message, error) {
	var messages []Message
	err := s.DB.Select(&messages, "SELECT id, message, user_id FROM messages WHERE conversation_id = $1", conversationId)
	if err != nil {
		return []Message{}, err
	}

	return messages, nil
}
