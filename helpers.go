package messengerbot

import (
	"github.com/satori/go.uuid"
)

func NewUserFromId(id string) User {
	return User{
		Id: id,
	}
}

func NewUserFromPhone(phoneNumber string) User {
	return User{
		PhoneNumber: phoneNumber,
	}
}

func NewMessage(text string) Message {
	return Message{
		Text: text,
	}
}

func NewImageMessage(url string) Message {
	return Message{
		Attachment: &Attachment{
			Type: "image",
			Payload: ImagePayload{
				Url: url,
			},
		},
	}
}

func NewGenericTemplate() GenericTemplate {
	return GenericTemplate{
		TemplateBase: TemplateBase{
			Type: "generic",
		},
		Elements: []Element{},
	}
}

func NewElement(title string) Element {
	return Element{
		Title: title,
	}
}

func NewWebUrlButton(title, url string) Button {
	return Button{
		Type:  "web_url",
		Title: title,
		Url:   url,
	}
}

func NewPostbackButton(title, postback string) Button {
	return Button{
		Type:    "postback",
		Title:   title,
		Payload: postback,
	}
}

func NewButtonTemplate(text string) ButtonTemplate {
	return ButtonTemplate{
		TemplateBase: TemplateBase{
			Type: "button",
		},
		Text:    text,
		Buttons: []Button{},
	}
}

func NewReceiptTemplate(recipientName string) ReceiptTemplate {
	return ReceiptTemplate{
		TemplateBase: TemplateBase{
			Type: "receipt",
		},
		RecipientName: recipientName,
		Id:            uuid.NewV4().String(),
		Currency:      "USD",
		PaymentMethod: "",
		Items:         []OrderItem{},
		Summary:       OrderSummary{},
	}
}
