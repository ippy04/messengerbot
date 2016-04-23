Facebook Messenger Platform (Chatbot) GoLang API
=====

A Golang implementation of the [Facebook Messenger Platform](https://developers.facebook.com/docs/messenger-platform).

## Installation
```bash
go get github.com/ippy04/messengerbot
```

Or if using the excellent [glide] (https://github.com/Masterminds/glide) package manager:

```bash
glide get github.com/ippy04/messengerbot
```

## Usage
```go

	bot := messengerbot.NewMessengerBot(accessToken, verifyToken)
	bot.Debug = true

	user := messengerbot.NewUserFromId(979665978736393)
	msg := messengerbot.NewMessage("Hello World")

	bot.Send(user, msg, messengerbot.NotificationTypeRegular)
```


