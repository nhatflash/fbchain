package helper

import (
	"time"
	appError "github.com/nhatflash/fbchain/error"
)

var dateFormat string = "2006-01-02"
var dateTimeFormat string = "2006-01-02 15:04:00"
var timeFormat string = "13:12:05"
var timeZone = "Asia/Bangkok"

func ConvertToDate(dateStr string) (*time.Time, error) {
	loc, locErr := time.LoadLocation(timeZone);
	if locErr != nil {
		return nil, appError.InternalError("Error of location: " + locErr.Error())
	}
	date, dateErr := time.ParseInLocation(dateFormat, dateStr, loc)
	if dateErr != nil {
		return nil, appError.InternalError("Error when parsing date: " + dateStr)
	}
	return &date, nil
}

func ConvertToDateTime(dateTimeStr string) (*time.Time, error) {
	loc, locErr := time.LoadLocation(timeZone);
	if locErr != nil {
		return nil, appError.InternalError("Error of location: " + locErr.Error())
	}
	dateTime, dateTimeErr := time.ParseInLocation(dateTimeFormat, dateTimeStr, loc)
	if dateTimeErr != nil {
		return nil, appError.InternalError("Error when parsing date time: " + dateTimeStr)
	}
	return &dateTime, nil
}


func ConvertToTime(timeStr string) (*time.Time, error) {
	loc, locErr := time.LoadLocation(timeZone);
	if locErr != nil {
		return nil, appError.InternalError("Error of location: " + locErr.Error())
	}
	time, timeErr := time.ParseInLocation(timeFormat, timeStr, loc)
	if timeErr != nil {
		return nil, appError.InternalError("Error when parsing time: " + timeStr)
	}
	return &time, nil
}