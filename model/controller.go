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

func PushBookings(booking Booking) {
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
	fmt.Println("rggflfglg", bookings)
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
			AvailableBy: common.UtcToIST(booking.ReservationEnd).String(),
			Blocked:     true,
		}
	}

	for id := range config.Rooms {
		_, exists := roomStatus[id]
		if exists == false {
			roomStatus[id] = RoomStatus{
				Blocked:  false,
				RoomName: config.Rooms[id],
			}
		}
	}
	return roomStatus
}

func GetAvailableRoom(reservationStart time.Time,
	reservationEnd time.Time,
	preferredRoom string) (bool, string) {

	bookings := getOverlap(reservationStart, reservationEnd)
	availableRooms := make(map[string]bool)

	if len(bookings) == len(config.Rooms) {
		return false, ""
	}

	if len(bookings) == 0 {
		return true, preferredRoom
	}

	for id, _ := range config.Rooms {
		availableRooms[id] = true
	}

	for _, booking := range bookings {
		availableRooms[booking.MeetingRoom] = false
	}

	if availableRooms[preferredRoom] {
		return true, preferredRoom
	}

	for room, isAvailable := range availableRooms {
		if isAvailable {
			return true, room
		}
	}

	return false, ""
}
