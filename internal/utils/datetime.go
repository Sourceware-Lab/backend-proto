package utils

import "time"

var DatetimeFormats = []string{ //nolint:gochecknoglobals
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}

func ParseAnyDatetime(s string) (time.Time, error) {
	parsedTime := time.Time{}

	for _, layout := range DatetimeFormats {
		if parsedTime, err := time.Parse(layout, s); err == nil {
			return parsedTime, nil
		}
	}

	return parsedTime, nil
}
