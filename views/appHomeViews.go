package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"time"

	"encoding/json"

	"github.com/akash/booker/model"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

//go:embed assets/*
var appHomeAssets embed.FS

func AppHomeTabView(bookings []model.Booking) slack.HomeTabViewRequest {

	str, err := appHomeAssets.ReadFile("assets/AppHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView`: %v", err)
	}
	view := slack.HomeTabViewRequest{}
	json.Unmarshal([]byte(str), &view)

	for _, booking := range bookings {
		t, err := template.ParseFS(appHomeAssets, "assets/BookingBlock.json")
		if err != nil {
			panic(err)
		}
		var tpl bytes.Buffer
		err = t.Execute(&tpl, booking)
		if err != nil {
			panic(err)
		}
		str, _ = ioutil.ReadAll(&tpl)
		block_view := slack.HomeTabViewRequest{}
		json.Unmarshal(str, &block_view)
		view.Blocks.BlockSet = append(view.Blocks.BlockSet, block_view.Blocks.BlockSet...)
	}
	return view
}

func CreateSchedulebookingModal() slack.ModalViewRequest {

	view := slack.ModalViewRequest{}
	t, err := template.ParseFS(appHomeAssets, "assets/scheduleBooking/ScheduleBookingModal.json")
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	now := time.Now()

	curr := map[string]string{
		"Defaultdate":      now.Format("2006-01-02"),
		"DefaultStartTime": now.Format("15:04"),
		"DefaultEndTime":   now.Add(time.Minute * 30).Format("15:04"),
	}
	err = t.Execute(&tpl, curr)
	if err != nil {
		panic(err)
	}
	str, _ := ioutil.ReadAll(&tpl)
	json.Unmarshal(str, &view)

	nv := slack.ModalViewRequest{}

	preferrenceOptions := generatePreferanceBlock()
	json.Unmarshal(preferrenceOptions, &nv)

	view.Blocks.BlockSet = append(view.Blocks.BlockSet, nv.Blocks.BlockSet...)
	return view
}

func CreateRoomStatusModal(rooms map[string]model.RoomStatus) slack.ModalViewRequest {
	view := slack.ModalViewRequest{}
	str, err := appHomeAssets.ReadFile("assets/status/StatusModal.json")
	if err != nil {
		log.Printf("Unable to read view `StatusModal`: %v", err)
	}
	json.Unmarshal(str, &view)

	for _, room := range rooms {
		if room.Blocked == false {
			t, err := template.ParseFS(appHomeAssets, "assets/status/RoomStatusAvailable.json")
			if err != nil {
				panic(err)
			}
			var tpl bytes.Buffer
			err = t.Execute(&tpl, room)
			if err != nil {
				panic(err)
			}
			str, _ = ioutil.ReadAll(&tpl)
			blockView := slack.ModalViewRequest{}
			json.Unmarshal(str, &blockView)

			view.Blocks.BlockSet = append(view.Blocks.BlockSet, blockView.Blocks.BlockSet...)
		} else {
			t, err := template.ParseFS(appHomeAssets, "assets/status/RoomStatusBlocked.json")
			if err != nil {
				panic(err)
			}
			var tpl bytes.Buffer
			err = t.Execute(&tpl, room)
			if err != nil {
				panic(err)
			}
			str, _ = ioutil.ReadAll(&tpl)
			blockView := slack.ModalViewRequest{}
			json.Unmarshal(str, &blockView)
			view.Blocks.BlockSet = append(view.Blocks.BlockSet, blockView.Blocks.BlockSet...)
		}
	}

	return view
}

func AppHomeCreateBookingSuccessLabel(bookings []model.Booking) slack.HomeTabViewRequest {
	str, err := appHomeAssets.ReadFile("assets/AppHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView`: %v", err)
	}
	view := slack.HomeTabViewRequest{}
	json.Unmarshal(str, &view)

	for _, booking := range bookings {
		t, err := template.ParseFS(appHomeAssets, "assets/BookingBlock.json")
		if err != nil {
			panic(err)
		}
		var tpl bytes.Buffer
		err = t.Execute(&tpl, booking)
		if err != nil {
			panic(err)
		}
		str, _ = ioutil.ReadAll(&tpl)
		block_view := slack.HomeTabViewRequest{}
		json.Unmarshal(str, &block_view)
		view.Blocks.BlockSet = append(view.Blocks.BlockSet, block_view.Blocks.BlockSet...)
	}
	return view
}


func ModalSubmissionResponse(callback slack.InteractionCallback, client *socketmode.Client) {
	str, _ := appHomeAssets.ReadFile("assets/status/StatusModal.json")
	view := slack.ModalViewRequest{}
	json.Unmarshal(str, &view)
	r, e := client.OpenView(callback.TriggerID, view)
	fmt.Println(r.Blocks, r.NotifyOnClose, "err=>", e)
}
