package model

import (
	"time"
)

func getUpcomingBookingByUserID(UserID string) ([]Booking, error) {
	var bookings []Booking
	timeLimit := time.Now().Add(time.Duration(-2) * time.Hour).UTC()
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
