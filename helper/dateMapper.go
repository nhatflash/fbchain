package helper

import (
	"time"
	appError "github.com/nhatflash/fbchain/error"
)

var dateFormat string = "2006-01-02"
var dateTimeFormat string = "2006-01-02 15:04:05"
var timeFormat string = "15:04:05"
var timeZone = "Asia/Bangkok"

func ConvertToDate(dateStr string) (*time.Time, error) {
	var err error
	var loc *time.Location
	loc, err = time.LoadLocation(timeZone);
	if err != nil {
		return nil, err
	}
	var date time.Time
	date, err = time.ParseInLocation(dateFormat, dateStr, loc)
	if err != nil {
		return nil, appError.BadRequestError("Error when parsing date: " + dateStr)
	}
	return &date, nil
}

func ConvertToDateTime(dateTimeStr string) (*time.Time, error) {
	var err error
	var loc *time.Location
	loc, err = time.LoadLocation(timeZone);
	if err != nil {
		return nil, err
	}
	var dateTime time.Time
	dateTime, err = time.ParseInLocation(dateTimeFormat, dateTimeStr, loc)
	if err != nil {
		return nil, appError.BadRequestError("Error when parsing date time: " + dateTimeStr)
	}
	return &dateTime, nil
}


func ConvertToTime(timeStr string) (*time.Time, error) {
	var err error
	var loc *time.Location
	loc, err = time.LoadLocation(timeZone);
	if err != nil {
		return nil, err
	}
	var timeF time.Time
	timeF, err = time.ParseInLocation(timeFormat, timeStr, loc)
	if err != nil {
		return nil, appError.BadRequestError("Error when parsing time: " + timeStr)
	}
	return &timeF, nil
}


func ConvertDateTimeToString(dateTime time.Time) string {
	return dateTime.Format(dateTimeFormat)
}


func ConvertDateToString(date time.Time) string {
	return date.Format(dateFormat)
}