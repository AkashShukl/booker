package views

import (
	"encoding/json"
	"fmt"

	"github.com/akash/booker/config"
)

type Option struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
}

type Text struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

type Element struct {
	Type        string   `json:"type"`
	Placeholder Text     `json:"placeholder"`
	Options     []Option `json:"options"`
	ActionID    string   `json:"action_id"`
}

type Block struct {
	Type    string  `json:"type"`
	BlockID string  `json:"block_id"`
	Element Element `json:"element"`
	Label   Text    `json:"label"`
}

type Modal struct {
	Title      Text    `json:"title"`
	Submit     Text    `json:"submit"`
	Type       string  `json:"type"`
	CallbackID string  `json:"callback_id"`
	Close      Text    `json:"close"`
	Blocks     []Block `json:"blocks"`
}

func generatePreferanceBlock() []byte {

	rooms := config.Rooms
	data, err := appHomeAssets.ReadFile("assets/scheduleBooking/PreferredOption.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var modal Modal
	err = json.Unmarshal(data, &modal)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}

	var modOptions []Option
	for id, name := range rooms {
		temp := fmt.Sprintf(`{
						"text": {
							"type": "plain_text",
							"text": "%s",
							"emoji": true
						},
						"value": "%s"
					}`, name, id)
		var opt Option
		err = json.Unmarshal([]byte(temp), &opt)
		if err != nil {
			fmt.Println("Unable to marshal option to temp", err.Error())
		}
		modOptions = append(modOptions, opt)
	}

	modal.Blocks[0].Element.Options = modOptions
	jsonStr, _ := json.Marshal(modal)
	fmt.Println(string(jsonStr))
	return (jsonStr)
}
