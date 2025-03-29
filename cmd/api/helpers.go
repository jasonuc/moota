package main

import (
	"time"
)

func parseDurationForConfig(input string, ptr *time.Duration) {
	duration, err := time.ParseDuration(input)
	if err != nil {
		return
	}

	(*ptr) = duration
}
