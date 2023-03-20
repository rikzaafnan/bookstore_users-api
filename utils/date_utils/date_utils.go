package date_utils

import "time"

const (
	apiDateLayout     = "02-01-2006T15:04:05Z"
	apiDateLayoutMini = "02-01-2006"
	apiDBLayout       = "2006-01-02 15:04:05"
)

func GetNowString() string {
	now := GetNow()

	return now.Format(apiDateLayout)
}

func GetNowDateMiniString() string {
	now := GetNow()

	return now.Format(apiDateLayoutMini)
}

func GetNow() time.Time {

	return time.Now().UTC()
}

func GetNowDBFormatString() string {

	return GetNow().Format(apiDBLayout)
}
