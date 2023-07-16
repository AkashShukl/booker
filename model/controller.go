package model

import (
	"fmt"
	"time"

	"github.com/akash/booker/common"
	"github.com/akash/booker/config"
)

type Booking struct {
	ID               uint `gorm:"primarykey"`
	UserID           string
	UserName         string
	ReservationStart time.Time
	ReservationEnd   time.Time
	MeetingRoom      string
	Active           bool `gorm:"default:true"`
}

type RoomStatus struct {
	RoomNO      string
	RoomName    string
	BlockedBy   string
	AvailableBy string
	Blocked     bool
}

func PushBooking(booking Booking) {
	createBooking(&booking)
}

func CancelBooking(bookingID string) error {
	_, err := cancelBooking(bookingID)
	if err != nil {
		fmt.Println("DEBUG: Error deleting booking", err.Error())
		return err
	}
	return nil
}

func GetBookings(userID string) []Booking {
	bookings, err := getUpcomingBookingByUserID(userID)
	if err != nil {
		fmt.Println("DEBUG: Error getting upcoming bookings", err.Error())
		return nil
	}
	for i, booking := range bookings {
		bookings[i].ReservationEnd = common.UtcToIST(booking.ReservationEnd)
		bookings[i].ReservationStart = common.UtcToIST(booking.ReservationStart)
	}

	return bookings
}

func GetRoomStatus() map[string]RoomStatus {
	blockedRooms := getBlockedRoom()
	fmt.Println(blockedRooms)

	// var roomStatus []RoomStatus
	roomStatus := make(map[string]RoomStatus)

	for _, booking := range blockedRooms {

		roomStatus[booking.MeetingRoom] = RoomStatus{
			RoomNO:      booking.MeetingRoom,
			RoomName:    config.Rooms[booking.MeetingRoom],
			BlockedBy:   booking.UserName,
			AvailableBy: common.UtcToIST(booking.ReservationEnd).Format("02 Jan 15:04"),
			Blocked:     true,
		}
	}

	for id := range config.Rooms {
		_, exists := roomStatus[id]
		if !exists {
			roomStatus[id] = RoomStatus{
				Blocked:  false,
				RoomName: config.Rooms[id],
			}
		}
	}
	return roomStatus
}

func GetAvailableRooms(reservationStart time.Time,
	reservationEnd time.Time) (bool, map[string]bool) {

	bookings := getOverlap(reservationStart, reservationEnd)
	availableRooms := make(map[string]bool)

	if len(bookings) == len(config.Rooms) {
		return false, nil
	}

	for id, _ := range config.Rooms {
		availableRooms[id] = true
	}

	for _, booking := range bookings {
		delete(availableRooms, booking.MeetingRoom)
	}

	return true, availableRooms
}
