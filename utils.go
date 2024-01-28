package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func relativeToAbsoluteTime(relativeTime string) (time.Time, error) {
	now := time.Now()

	var unit string
	var quantityStr string

	// Check if the relative time ends with "mon"
	if strings.HasSuffix(relativeTime, "mon") {
		unit = "mon"
		quantityStr = relativeTime[:len(relativeTime)-3]
	} else {
		// Get the unit of time and the quantity
		unit = relativeTime[len(relativeTime)-1:]
		quantityStr = relativeTime[:len(relativeTime)-1]
	}

	// Convert the quantity to an integer
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return time.Time{}, err
	}

	// Subtract the quantity of the unit of time from the current time
	switch strings.ToLower(unit) {
	case "d":
		return now.AddDate(0, 0, -quantity), nil
	case "h":
		return now.Add(-time.Duration(quantity) * time.Hour), nil
	case "m":
		return now.Add(-time.Duration(quantity) * time.Minute), nil
	case "w":
		return now.AddDate(0, 0, -7*quantity), nil
	case "mon":
		return now.AddDate(0, -quantity, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown time unit: %s", unit)
	}
}
