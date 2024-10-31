package models

import "github.com/gorilla/websocket"

type WebsocketConnection struct {
	ID   string
	Type string
	Conn *websocket.Conn
}

type BondedConnection struct {
	ConnUser    *WebsocketConnection
	ConnSupport *WebsocketConnection
	ID          string
}
