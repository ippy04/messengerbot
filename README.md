Facebook Messenger Platform (Chatbot) Go API
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

### Send a Regular Message
```go

	bot := messengerbot.NewMessengerBot(accessToken, verifyToken)
	bot.Debug = true
	
	user := messengerbot.NewUserFromId(userId)
	msg := messengerbot.NewMessage("Hello World")

	bot.Send(user, msg, messengerbot.NotificationTypeRegular)
```
### Send an Image Message
```go
	msg := messengerbot.NewImageMessage("https://pixabay.com/static/uploads/photo/2016/04/01/09/29/cartoon-1299393_960_720.png")
	bot.Send(user, msg, messengerbot.NotificationTypeRegular)
```


### Send an [Button Template](https://developers.facebook.com/docs/messenger-platform/send-api-reference#button_template) Message
```go
	msg := messengerbot.NewButtonTemplate("Pick one, any one")
	button1 := messengerbot.NewPostbackButton("Button 1", "POSTBACK_BUTTON_1")
	button2 := messengerbot.NewPostbackButton("Button 2", "POSTBACK_BUTTON_2")
	button3 := messengerbot.NewPostbackButton("Button 3", "POSTBACK_BUTTON_3")
	msg.AddButton(button1, button2, button3)
	
	bot.Send(user, msg, messengerbot.NotificationTypeRegular)
```


### Send a [Generic Template](https://developers.facebook.com/docs/messenger-platform/send-api-reference#generic_template) Message
```go
	msg := messengerbot.NewGenericTemplate()
	element := messengerbot.Element{
		Title:    "This is a bolded title",
		ImageUrl: "https://pixabay.com/static/uploads/photo/2016/04/01/09/29/cartoon-1299393_960_720.png",
		Subtitle: "I am a dinosaur. Hear me Rawr.",
	}
	
	button1 := messengerbot.NewPostbackButton("Button 1", "POSTBACK_BUTTON_1")
	button2 := messengerbot.NewPostbackButton("Button 2", "POSTBACK_BUTTON_2")
	button3 := messengerbot.NewPostbackButton("Button 3", "POSTBACK_BUTTON_3")
	element.AddButton(button1, button2, button3)
	
	msg.AddElement(element)
	
	bot.Send(user, msg, messengerbot.NotificationTypeRegular)
```


### Push Notification Types
```go
bot.Send(user, msg, messengerbot.NotificationTypeRegular)     // regular sound, vibrate and phone alert
bot.Send(user, msg, messengerbot.NotificationTypeSilentPush)  // phone notification only, no sound or vibrate alert
bot.Send(user, msg, messengerbot.NotificationTypeNoPush)      // no sound or phone notification
```
