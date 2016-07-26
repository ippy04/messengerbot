package messengerbot

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type rawEvent struct {
	Object  string          `json:"object"`
	Entries []*MessageEvent `json:"entry"`
}

type Event struct {
	Id   string `json:"id"`
	Time int64 `json:"time"`
}

type MessageEvent struct {
	Event
	Messaging []struct {
		MessageOpts
		Message  *ReceivedMessage `json:"message,omitempty"`
		Delivery *Delivery        `json:"delivery,omitempty"`
		Postback *Postback        `json:"postback,omitempty"`
		Optin    *Optin           `json:"optin,empty"`
	} `json:"messaging"`
}

type ReceivedMessage struct {
	Message
	Id  string `json:"mid"`
	Seq int    `json:"seq"`
}

type Delivery struct {
	MessageIds []string `json:"mids"`
	Watermark  int64    `json:"watermark"`
	Seq        int      `json:"seq"`
}

type Postback struct {
	Payload string `json:"payload"`
}

type Optin struct {
	Ref string `json:"ref"`
}

type MessageOpts struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`

	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`

	Timestamp int64 `json:"timestamp"`
}

type MessageReceivedHandler func(Event, MessageOpts, ReceivedMessage, chan bool)
type MessageDeliveredHandler func(Event, MessageOpts, Delivery, chan bool)
type PostbackHandler func(Event, MessageOpts, Postback, chan bool)
type AuthenticationHandler func(Event, MessageOpts, *Optin, chan bool)
type ErrorHandler func(Event, MessageOpts)

// Handler is the main HTTP handler for the Messenger service.
// It must be attached to some web server in order to receive messages
func (bot *MessengerBot) Handler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		query := req.URL.Query()
		if query.Get("hub.verify_token") != bot.VerifyToken {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(query.Get("hub.challenge")))
	} else if req.Method == "POST" {
		bot.handlePOST(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (bot *MessengerBot) handlePOST(rw http.ResponseWriter, req *http.Request) {
	read, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	//Message integrity check
	if bot.AppSecret != "" {
		if !checkIntegrity(bot.AppSecret, read, req.Header.Get("x-hub-signature")[5:]) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	event := &rawEvent{}
	err = json.Unmarshal(read, event)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range event.Entries {
		for _, message := range entry.Messaging {
			if message.Delivery != nil {
				if bot.MessageDelivered != nil {
					chHasError := make(chan bool)
					go bot.MessageDelivered(entry.Event, message.MessageOpts, *message.Delivery, chHasError)
					bot.handleError(entry.Event, message.MessageOpts, chHasError)
				}
			} else if message.Message != nil {
				if bot.MessageReceived != nil {
					chHasError := make(chan bool)
					go bot.MessageReceived(entry.Event, message.MessageOpts, *message.Message, chHasError)
					bot.handleError(entry.Event, message.MessageOpts, chHasError)
				}
			} else if message.Postback != nil {
				if bot.Postback != nil {
					chHasError := make(chan bool)
					go bot.Postback(entry.Event, message.MessageOpts, *message.Postback, chHasError)
					bot.handleError(entry.Event, message.MessageOpts, chHasError)
				}
			} else if bot.Authentication != nil {
				chHasError := make(chan bool)
				go bot.Authentication(entry.Event, message.MessageOpts, message.Optin, chHasError)
				bot.handleError(entry.Event, message.MessageOpts, chHasError)
			}
		}
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

func (bot *MessengerBot) handleError(event Event, opts MessageOpts, chHasError chan bool) {
	hasError := <-chHasError
	close(chHasError)
	if hasError {
		go bot.Error(event, opts)
	}
}

func checkIntegrity(appSecret string, bytes []byte, expectedSignature string) bool {
	mac := hmac.New(sha1.New, []byte(appSecret))
	if fmt.Sprintf("%x", mac.Sum(bytes)) != expectedSignature {
		return false
	}
	return true
}
