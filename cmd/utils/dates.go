package utils

import "time"

func IsBeforeToday(dateStr string, now time.Time) bool {
	date, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		return false
	}

	today := now.Truncate(24 * time.Hour)
	dateWithoutTime := date.Truncate(24 * time.Hour)

	return dateWithoutTime.Before(today)
}

func IsToday(dateStr string, now time.Time) bool {
	date, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		panic(err)
	}

	today := now.Truncate(24 * time.Hour)
	dateWithoutTimezone := date.UTC().Truncate(24 * time.Hour)

	return dateWithoutTimezone.Equal(today)
}

func GetDateOffsetValue(date time.Time) float64 {
	_, zone := date.Zone()
	offset := (float64)(zone) / 60 / 60

	return offset
}

func AddOffsetToDate(date time.Time, offset float64) time.Time {
	timeToAdd := offset * -1
	modifiedDate := date.Add(time.Hour * time.Duration(timeToAdd))

	return modifiedDate
}
