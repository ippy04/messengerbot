package messengerbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
)

type ctaBase struct {
	SettingType string `json:"setting_type"`
	ThreadState string `json:"thread_state"`
}

var welcomeMessage = ctaBase{
	SettingType: "setting_type",
	ThreadState: "new_thread",
}

type Settings struct {
	SettingType string `json:"setting_type"`
	Greeting Greeting `json:"greeting"`
}

type Greeting struct {
	Text string `json:"text"`
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

// SetWelcomeMessage sets the message that is sent first. If message is nil or empty the welcome message is not sent.
func (bot *MessengerBot) SetGreeting() error {
	setting := Settings{
			SettingType: "greeting",
			Greeting: Greeting{
				Text:"Welcome to the bots",
			},
	}

	byt, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	fmt.Print(string(byt))

	url := GraphAPI+"me/thread_settings?access_token="+bot.AccessToken
	request, err := http.NewRequest("POST", url, bytes.NewReader(byt))
	fmt.Print(url)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	response, err := bot.Client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Invalid status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
  fmt.Println("response Body:", string(body))

	return nil
}

// SetWelcomeMessage sets the message that is sent first. If message is nil or empty the welcome message is not sent.
func (bot *MessengerBot) SetGetstarted() error {
	str := `
	{
	  "setting_type":"call_to_actions",
	  "thread_state":"new_thread",
	  "call_to_actions":[
	    {
	      "payload":"GetStarted"
	    }
	  ]
	}
	`

	byt := []byte(str)
	fmt.Print(string(byt))

	url := GraphAPI+"me/thread_settings?access_token="+bot.AccessToken
	request, err := http.NewRequest("POST", url, bytes.NewReader(byt))
	fmt.Print(url)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	response, err := bot.Client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Invalid status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
  fmt.Println("response Body:", string(body))

	return nil
}
