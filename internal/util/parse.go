package util

import (
	"fmt"
	"strconv"
	"time"
)

// ParseAgeString parses strings like "30d", "2m", "1y" into a time.Duration
func ParseAgeString(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid age string: %s", s)
	}
	unit := s[len(s)-1]
	numStr := s[:len(s)-1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number in age string: %s", s)
	}
	switch unit {
	case 'd':
		return time.Duration(num) * 24 * time.Hour, nil
	case 'm':
		return time.Duration(num) * 30 * 24 * time.Hour, nil // approx month
	case 'y':
		return time.Duration(num) * 365 * 24 * time.Hour, nil // approx year
	default:
		return 0, fmt.Errorf("invalid unit in age string: %s", s)
	}
}
