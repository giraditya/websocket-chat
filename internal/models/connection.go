package models

import "github.com/gorilla/websocket"

type WebsocketConnection struct {
	Type string
	Conn *websocket.Conn
}

type BondedConnection struct {
	ConnUser    *WebsocketConnection
	ConnSupport *WebsocketConnection
	ChatID      string
}
