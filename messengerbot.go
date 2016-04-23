package messengerbot

import (
	"net/http"
)

// Messenger is the main service which handles all callbacks from facebook
// Events are delivered to handlers if they are specified
type MessengerBot struct {
	MessageReceived  MessageReceivedHandler
	MessageDelivered MessageDeliveredHandler
	Postback         PostbackHandler
	Authentication   AuthenticationHandler

	VerifyToken string
	AppSecret   string // Optional: For validating integrity of messages
	AccessToken string
	PageId      string // Optional: For setting welcome message
	Debug       bool
	Client      *http.Client
}

func NewMessengerBot(token string, verifyToken string) *MessengerBot {
	return &MessengerBot{
		AccessToken: token,
		VerifyToken: verifyToken,
		Debug:       false,
		Client:      &http.Client{},
	}
}
