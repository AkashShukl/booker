package handlers

import (
	"fmt"
	"log"

	"github.com/akash/booker/common"
	"github.com/akash/booker/model"
	"github.com/akash/booker/views"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func AppHomeOpenedEvent(ev *slackevents.AppHomeOpenedEvent, client *socketmode.Client) {
	if ev.Tab != "home" {
		return
	}
	user := ev.User
	println("user =>  ", user)
	view := views.AppHomeTabView(model.GetBookings(user))
	_, err := client.PublishView(user, view, "")
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func AppHomeDeclineBooking(bookingID string, callback slack.InteractionCallback, client *socketmode.Client) {
	user := callback.User.ID
	println("user =>  ", user)
	model.CancelBooking(bookingID)
	view := views.AppHomeTabView(model.GetBookings(user))
	_, err := client.PublishView(user, view, "")
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func PublishScheduleBooking(callback slack.InteractionCallback, client *socketmode.Client) {

	date := callback.View.State.Values[views.ModalScheduleDateBlockID]["datepicker_action"].SelectedDate
	startTime := callback.View.State.Values[views.ModalScheduleTimeStartBlockID]["timepicker_action_from"].SelectedTime
	endTime := callback.View.State.Values[views.ModalScheduleTimeEndBlockID]["timepicker_action_to"].SelectedTime

	reservationStart := common.DateTimeToUTC(date+" "+startTime, "2006-01-02 15:04")
	reservationEnd := common.DateTimeToUTC(date+" "+endTime, "2006-01-02 15:04")

	preferredRoom := callback.View.State.Values["room_preference"]["static_select-action"].SelectedOption.Value
	// check for room availability
	ok, room := model.GetAvailableRoom(reservationStart, reservationEnd, preferredRoom)

	if ok == false {
		fmt.Println("Sorry No rooms available!")
	} else {
		booking := model.Booking{
			UserID:           callback.User.ID,
			UserName:         callback.User.Name,
			ReservationStart: reservationStart,
			ReservationEnd:   reservationEnd,
			MeetingRoom:      room,
			Active:           true,
		}
		model.PushBookings(booking)
		view := views.AppHomeCreateBookingSuccessLabel(model.GetBookings(callback.User.ID))
		_, err := client.PublishView(callback.User.ID, view, "")
		if err != nil {
			log.Printf("ERROR createStickieNote: %v", err)
		}
	}

}
