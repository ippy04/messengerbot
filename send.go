package messengerbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NotificationType string

const (
	NotificationTypeRegular    NotificationType = "REGULAR"     // regular sound, vibrate and phone alert
	NotificationTypeSilentPush NotificationType = "SILENT_PUSH" // phone notification only, no sound or vibrate alert
	NotificationTypeNoPush     NotificationType = "NO_PUSH"     // no sound or phone notification
)

const (
	GraphAPI             = "https://graph.facebook.com/v2.6/"
	MessengerAPIEndpoint = GraphAPI + "me/messages?access_token=%s"
	ProfileAPIEndpoint   = GraphAPI + "%d?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%s"

	GenericTemplateTitleLengthLimit       = 45
	GenericTemplateSubtitleLengthLimit    = 80
	GenericTemplateCallToActionTitleLimit = 20
	GenericTemplateCallToActionItemsLimit = 3
	GenericTemplateBubblesPerMessageLimit = 10

	ButtonTemplateButtonsLimit = 3
)

var (
	ErrTitleLengthExceeded             = errors.New("Template element title exceeds the 45 character limit")
	ErrSubtitleLengthExceeded          = errors.New("Template element subtitle exceeds the 80 character limit")
	ErrCallToActionTitleLengthExceeded = errors.New("Template call to action title exceeds the 20 character limit")
	ErrButtonsLimitExceeded            = errors.New("Limit of 3 buttons exceeded")
	ErrBubblesLimitExceeded            = errors.New("Limit of 10 bubbles per message exceeded")
)

type User struct {
	Id          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string            `json:"type"`
	Payload AttachmentPayload `json:"payload"`
}

type AttachmentPayload interface{}

type ImagePayload struct {
	Url string `json:"url"`
}

type TemplateBase struct {
	Type string `json:"template_type"`
}

type GenericTemplate struct {
	TemplateBase
	Elements []Element `json:"elements"`
}

func (g GenericTemplate) Validate() error {
	if len(g.Elements) > GenericTemplateBubblesPerMessageLimit {
		return ErrBubblesLimitExceeded
	}
	for _, elem := range g.Elements {
		if len(elem.Title) > GenericTemplateTitleLengthLimit {
			return ErrTitleLengthExceeded
		}

		if len(elem.Subtitle) > GenericTemplateSubtitleLengthLimit {
			return ErrSubtitleLengthExceeded
		}

		if len(elem.Buttons) > GenericTemplateCallToActionItemsLimit {
			return ErrButtonsLimitExceeded
		}

		for _, button := range elem.Buttons {
			if len(button.Title) > GenericTemplateCallToActionTitleLimit {
				return ErrCallToActionTitleLengthExceeded
			}
		}
	}
	return nil
}

func (g *GenericTemplate) AddElement(e ...Element) {
	g.Elements = append(g.Elements, e...)
}

type Element struct {
	Title    string   `json:"title"`
	Url      string   `json:"item_url,omitempty"`
	ImageUrl string   `json:"image_url,omitempty"`
	Subtitle string   `json:"subtitle,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`
}

func (e *Element) AddButton(b ...Button) {
	e.Buttons = append(e.Buttons, b...)
}

type Button struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Url     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type ButtonTemplate struct {
	TemplateBase
	Text    string   `json:"text,omitempty"`
	Buttons []Button `json:"buttons,omitempty"`
}

func (b ButtonTemplate) Validate() error {
	if len(b.Buttons) > ButtonTemplateButtonsLimit {
		return ErrButtonsLimitExceeded
	}
	return nil
}

func (b *ButtonTemplate) AddButton(bt ...Button) {
	b.Buttons = append(b.Buttons, bt...)
}

type ReceiptTemplate struct {
	TemplateBase
	RecipientName string            `json:"recipient_name"`
	Id            string            `json:"order_number"`
	Currency      string            `json:"currency"`
	PaymentMethod string            `json:"payment_method"`
	Timestamp     int64             `json:"timestamp,omitempty"`
	Url           string            `json:"order_url,omitempty"`
	Items         []OrderItem       `json:"elements"`
	Address       *OrderAddress     `json:"address,omitempty"`
	Summary       OrderSummary      `json:"summary"`
	Adjustments   []OrderAdjustment `json:"adjustments,omitempty"`
}

type OrderItem struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Quantity int64  `json:"quantity,omitempty"`
	Price    int64  `json:"price,omitempty"`
	Currency string `json:"currency,omitempty"`
	ImageURL string `json:"image_url,omiempty"`
}

type OrderAddress struct {
	Street1    string `json:"street_1"`
	Street2    string `json:"street_2,omitempty"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

type OrderSummary struct {
	TotalCost    int `json:"total_cost,omitempty"`
	Subtotal     int `json:"subtotal,omitempty"`
	ShippingCost int `json:"shipping_cost,omitempty"`
	TotalTax     int `json:"total_tax,omitempty"`
}

type OrderAdjustment struct {
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
}

type SendRequest struct {
	Recipient        User    `json:"recipient"`
	Message          Message `json:"message"`
	NotificationType string  `json:"notification_type"`
}

type SendResponse struct {
	RecipientId string        `json:"recipient_id"`
	MessageId   string        `json:"message_id"`
	Error       ErrorResponse `json:"error"`
}
type ErrorResponse struct {
	Message    string `json:"message"`
	Type       string `json:"type"`
	Code       int64  `json:"code"`
	ErrorData  string `json:"error_data"`
	FBstraceId string `json:"fbstrace_id"`
}

func (bot *MessengerBot) MakeRequest(byt *bytes.Buffer) (*SendResponse, error) {
	url := fmt.Sprintf(MessengerAPIEndpoint, bot.AccessToken)

	request, _ := http.NewRequest("POST", url, byt)
	request.Header.Set("Content-Type", "application/json")
	response, err := bot.Client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(body, &er)

		if bot.Debug {
			log.Printf("[DEBUG] Response: %s", string(body))
		}

		return nil, errors.New("Error response received: " + er.Error.Message)
	}

	var sendResponse *SendResponse
	json.Unmarshal(body, &sendResponse)
	if err != nil {
		return nil, err
	}

	return sendResponse, nil
}

// convenience method for sending a text-based message
func (bot *MessengerBot) SendTextMessage(user User, messageText string, notificationType NotificationType) (*SendResponse, error) {
	message := NewMessage(messageText)
	return bot.Send(user, message, notificationType)
}

func (bot *MessengerBot) Send(user User, content interface{}, notificationType NotificationType) (*SendResponse, error) {
	var r SendRequest

	switch content.(type) {
	case Message:
		r = SendRequest{
			Recipient:        user,
			Message:          content.(Message),
			NotificationType: string(notificationType),
		}

	case GenericTemplate:
		r = SendRequest{
			Recipient:        user,
			NotificationType: string(notificationType),
			Message: Message{
				Attachment: &Attachment{
					Type:    "template",
					Payload: content.(GenericTemplate),
				},
			},
		}

	case ButtonTemplate:
		r = SendRequest{
			Recipient:        user,
			NotificationType: string(notificationType),
			Message: Message{
				Attachment: &Attachment{
					Type:    "template",
					Payload: content.(ButtonTemplate),
				},
			},
		}

	case ReceiptTemplate:
		r = SendRequest{
			Recipient:        user,
			NotificationType: string(notificationType),
			Message: Message{
				Attachment: &Attachment{
					Type:    "template",
					Payload: content.(ReceiptTemplate),
				},
			},
		}

	default:
		return nil, errors.New("Unsupported message content type")
	}

	if r == (SendRequest{}) {
		return nil, errors.New("Unknown Error - Unable to create SendRequest")
	}

	payload, _ := json.Marshal(r)
	if bot.Debug {
		log.Printf("[DEBUG] Payload: %s", string(payload))
	}
	return bot.MakeRequest(bytes.NewBuffer(payload))
}
