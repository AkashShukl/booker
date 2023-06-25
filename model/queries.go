package model

import (
	"fmt"
	"time"
)

// All the DB operations should be moved to different file

func getAllBookingByUserIDDebug(UserID string) {
	fmt.Println("ALL BOOKINGS ")
	var bookings []Booking
	_ = db.Where("user_id = ?  AND active = true",
		UserID).Find(&bookings).Error

	for i, book := range bookings {
		fmt.Println(i, "=> ", book)
	}

}

func getUpcomingBookingByUserID(UserID string) ([]Booking, error) {

	getAllBookingByUserIDDebug(UserID)
	var bookings []Booking
	timeLimit := time.Now().Add(time.Duration(-2) * time.Hour).UTC()
	fmt.Println("timeLimit: ", timeLimit)
	err := db.Where("user_id = ? and reservation_start > ? AND active = true",
		UserID, timeLimit).Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil

}

func createBooking(booking *Booking) (bool, error) {
	err := db.Create(&booking).Error
	if err != nil {
		return false, err
	}
	fmt.Println("DEBUG: Booking successfully created")
	return true, nil
}

func cancelBooking(bookingID string) (bool, error) {
	var booking Booking

	err := db.First(&booking, bookingID).Error
	if err != nil {
		panic("failed to retrieve booking")
	}

	booking.Active = false
	err = db.Save(&booking).Error
	if err != nil {
		panic("failed to update booking")
	}

	return true, nil
}

func getOverlap(searchStart, searchEnd time.Time) []Booking {
	var bookings []Booking

	err := db.Where("reservation_start BETWEEN ? AND ? OR reservation_end BETWEEN ? AND ? And active = true",
		searchStart, searchEnd, searchStart, searchEnd).
		Find(&bookings).Error
	if err != nil {
		panic("failed to retrieve data")
	}
	fmt.Println("Overlapping time => ")
	for _, booking := range bookings {
		fmt.Printf("ID: %d, Start: %s, End: %s\n",
			booking.ID,
			booking.ReservationStart,
			booking.ReservationEnd)
	}
	return bookings
}

func getBlockedRoom() []Booking {
	var bookings []Booking
	now := time.Now().UTC()
	err := db.Where("reservation_start < ? AND reservation_end > ? And active = true", now, now).Find(&bookings).Error
	if err != nil {
		panic("failed to retrieve data")
	}
	return bookings
}
