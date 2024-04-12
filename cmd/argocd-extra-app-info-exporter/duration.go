package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DurationFlag time.Duration

func (d *DurationFlag) Set(value string) error {
	// Extracting the numerical part of the input.
	val, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return err
	}

	// Determining the unit from the input.
	unit := value[len(value)-1]
	switch strings.ToLower(string(unit)) {
	case "s":
		*d = DurationFlag(time.Duration(val) * time.Second)
	case "m":
		*d = DurationFlag(time.Duration(val) * time.Minute)
	case "h":
		*d = DurationFlag(time.Duration(val) * time.Hour)
	default:
		return fmt.Errorf("invalid unit: %s", unit)
	}
	return nil
}

func (d *DurationFlag) String() string {
	return time.Duration(*d).String()
}
