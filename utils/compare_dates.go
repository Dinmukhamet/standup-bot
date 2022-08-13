package utils

import (
	"time"

	"github.com/Dinmukhamet/gostandup/constants"
)

func CountDaysBetweenTwoDates(date1, date2 string) int {
	t1, _ := time.Parse(constants.DATE_FORMAT, date1)
	t2, _ := time.Parse(constants.DATE_FORMAT, date2)
	return int(t2.Sub(t1).Hours() / 24)
}
