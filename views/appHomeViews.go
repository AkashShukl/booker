package views

import (
	"bytes"
	"embed"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"encoding/json"

	"github.com/akash/booker/config"
	"github.com/akash/booker/model"
	"github.com/slack-go/slack"
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

func AppHomeCreateBookingLabel(bookings []model.Booking,
	success bool,
	message string,
	availableRooms map[string]bool) slack.HomeTabViewRequest {

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

	if success {
		return view
	}

	t, err := template.ParseFS(appHomeAssets, "assets/PreferredRoomNABlock.json")
	if err != nil {
		panic(err)
	}
	msg := make(map[string]string)
	msg["Msg"] = message
	if availableRooms != nil {
		var str string
		for k, _ := range availableRooms {
			str += config.Rooms[k] + ", "
		}
		str = strings.TrimSuffix(str, ", ")
		msg["Options"] = str
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, msg)
	if err != nil {
		panic(err)
	}
	str, _ = ioutil.ReadAll(&tpl)
	faliureBlock := slack.HomeTabViewRequest{}
	json.Unmarshal(str, &faliureBlock)
	view.Blocks.BlockSet = append(view.Blocks.BlockSet, faliureBlock.Blocks.BlockSet...)
	return view

}
