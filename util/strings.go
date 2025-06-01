package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for i := 1; i <= tmpCount; i++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(i), 1)
	}
	return old
}

func ConvertToDate(date string, sep string) time.Time {
	timeRemoved := strings.Split(date, "T")[0]
	dateArray := strings.Split(timeRemoved, sep)
	var year, month, day int
	var err error
	if len(dateArray[0]) == 4 {
		year, err = strconv.Atoi(dateArray[0])
		if err != nil {
			fmt.Println(err)
		}
		month, err = strconv.Atoi(dateArray[1])
		if err != nil {
			fmt.Println(err)
		}
		day, err = strconv.Atoi(dateArray[2])
		if err != nil {
			fmt.Println(err)
		}
	} else {
		year, err = strconv.Atoi(dateArray[2])
		if err != nil {
			fmt.Println(err)
		}
		month, err = strconv.Atoi(dateArray[1])
		if err != nil {
			fmt.Println(err)
		}
		day, err = strconv.Atoi(dateArray[0])
		if err != nil {
			fmt.Println(err)
		}
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func FormatTime(t string) string {
	timeArray := strings.Split(strings.Split(t, "T")[1], ":")
	return timeArray[0] + ":" + timeArray[1]
}
