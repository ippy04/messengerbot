package messengerbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Profile struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_pic,omitempty"`
	Locale         string `json:"locale,omitempty"`
	Timezone       int    `json:"timezone,omitempty"`
	Gender         string `json:"gender,omitempty"`
}

func (bot *MessengerBot) GetProfile(userID string) (*Profile, error) {
	resp, err := bot.Client.Get(fmt.Sprintf(ProfileAPIEndpoint, userID, bot.AccessToken))

	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Body)

	read, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured fetching profile: " + er.Error.Message)
	}
	profile := new(Profile)
	return profile, json.Unmarshal(read, profile)
}
