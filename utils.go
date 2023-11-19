package main

import (
	"os"
	"regexp"
	"strconv"
	"time"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func str2time(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func IReplace(subject string, search string, replace string) string {
	searchRegex := regexp.MustCompile("(?i)" + search + "[\\s:]+")
	return searchRegex.ReplaceAllString(subject, replace)
}
