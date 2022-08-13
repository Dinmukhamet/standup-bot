package constants

import "time"

const (
	DATE_FORMAT           string = "2006-01-02"
	TIME_FORMAT           string = "03:04 PM"
	TIME_CRON_FORMAT      string = "03:04"
	DEFAULT_ERROR_MESSAGE string = "Что-то пошло не так 😔"
	DEADLINE_TIME         string = "04:40 AM"
	DEFAULT_TIMEZONE      string = "Asia/Bishkek"
)

var LOCATION *time.Location = nil
