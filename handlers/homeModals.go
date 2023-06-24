package handlers

import (
	"fmt"
	"log"

	"github.com/akash/booker/model"
	"github.com/akash/booker/views"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func AppHomeScheduledBookingModal(triggerID string, client *socketmode.Client) {
	view := views.CreateSchedulebookingModal()
	_, err := client.OpenView(triggerID, view)
	if err != nil {
		log.Printf("ERROR openCreateSchedulebookingModal: %v", err)
	}
}

func PublishCurrentRoomStatus(callback slack.InteractionCallback, client *socketmode.Client) {
	fmt.Println("triggered: PublishCurrentRoomStatus")
	status := model.GetRoomStatus()
	triggerID := callback.TriggerID
	view := views.CreateRoomStatusModal(status)
	_, err := client.OpenView(triggerID, view)
	if err != nil {
		log.Printf("ERROR openStatsModal: %v", err)
	}
}
