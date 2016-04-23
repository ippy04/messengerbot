package messengerbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ctaBase struct {
	SettingType string `json:"setting_type"`
	ThreadState string `json:"thread_state"`
}

var welcomeMessage = ctaBase{
	SettingType: "setting_type",
	ThreadState: "new_thread",
}

type cta struct {
	ctaBase
	CallToActions []*Message `json:"call_to_actions"`
}

type result struct {
	Result string `json:"result"`
}

// SetWelcomeMessage sets the message that is sent first. If message is nil or empty the welcome message is not sent.
func (bot *MessengerBot) SetWelcomeMessage(message *Message) error {
	cta := &cta{
		ctaBase:       welcomeMessage,
		CallToActions: []*Message{message},
	}
	if bot.PageId == "" {
		return errors.New("PageId is empty")
	}
	byt, err := json.Marshal(cta)
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", fmt.Sprintf(GraphAPI+"%s/thread_settings", bot.PageId), bytes.NewReader(byt))
	request.Header.Set("Content-Type", "application/json")

	response, err := bot.Client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Invalid status code")
	}
	decoder := json.NewDecoder(response.Body)
	result := &result{}
	err = decoder.Decode(result)
	if err != nil {
		return err
	}
	if result.Result != "Successfully added new_thread's CTAs" {
		return errors.New("Something went wrong with setting thread's welcome message, facebook error: " + result.Result)
	}
	return nil
}
