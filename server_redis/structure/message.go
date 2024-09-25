package structure

import (
	"time"
)

type RoomList struct {
	Type  string   `json:"type"`
	Rooms []string `json:"rooms"`
}

type UserList struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}

const (
	RequestJoin    string = "JOIN"
	RequestMessage string = "MESSAGE"
)

type Request struct {
	Method  string   `json:"method"  validate:"required,oneof=JOIN MESSAGE"`
	Name    *string  `json:"name"    validate:"required_if=Method JOIN,omitempty,gt=0,lt=13,alphanumeric"`
	Message *Message `json:"message" validate:"required_if=Method MESSAGE,omitempty"`
	Nonce   string   `json:"nonce"   validate:"required,uuid4"`
}

type Message struct {
	Room string `json:"room" validate:"required,gt=0,lt=11,alphanumeric"`
	Text string `json:"text" validate:"required,gt=0,lt=300"`
}

const (
	NotificationResponse string = "RESPONSE"
	NotificationMessage  string = "MESSAGE"
	NotificationJoin     string = "JOIN"
	NotificationLeave    string = "LEAVE"
)

type Notification struct {
	Type    string   `json:"type"`
	Content *Content `json:"content"`
	Nonce   string   `json:"nonce"`
	Error   string   `json:"error"` // "" = success
}

type Content struct {
	From     string     `json:"from"`
	Message  *Message   `json:"message"`
	DateTime *time.Time `json:"datetime"`
}

const (
	Success     string = ""
	NameError   string = "duplicate name"
	RoomError   string = "no room found"
	MethodError string = "method not valid"
)
