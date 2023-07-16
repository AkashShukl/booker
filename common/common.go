package common

import (
	"fmt"
	"time"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("Asia/Calcutta")
}

func ExtractTimeString(utcTime time.Time) string {
	localTime := utcTime.In(loc)
	hour := localTime.Hour()
	minute := localTime.Minute()
	formattedTime := fmt.Sprintf("%02d:%02d", hour, minute)
	return formattedTime
}

func ExtractDateString(utcTime time.Time) string {
	localTime := utcTime.In(loc)
	year := localTime.Year()
	month := localTime.Month()
	day := localTime.Day()
	formattedDate := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return formattedDate
}

func DateTimeToUTC(dateString, layout string) time.Time {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	loc, _ := time.LoadLocation("Asia/Calcutta")
	sourceTime, _ := time.ParseInLocation(layout, dateString, loc)
	utcTime := sourceTime.UTC()
	return utcTime
}

func UtcToIST(utcTime time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Calcutta")
	return utcTime.In(loc)
}

