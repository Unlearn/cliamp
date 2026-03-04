package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseJumpTarget parses an absolute jump target from track start.
//
// Accepted formats:
// - "10"      => 10s
// - "58:05"   => 58m05s
// - "58:5"    => 58m05s
// - "58:"     => 58m00s
func parseJumpTarget(raw string) (time.Duration, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("enter time like 10 or 58:05")
	}

	if strings.Count(s, ":") == 0 {
		if !isDigits(s) {
			return 0, fmt.Errorf("seconds must be a number")
		}
		secs, err := strconv.Atoi(s)
		if err != nil || secs < 0 {
			return 0, fmt.Errorf("invalid seconds value")
		}
		return time.Duration(secs) * time.Second, nil
	}

	if strings.Count(s, ":") != 1 {
		return 0, fmt.Errorf("use mm:ss format")
	}

	parts := strings.SplitN(s, ":", 2)
	minPart := strings.TrimSpace(parts[0])
	secPart := strings.TrimSpace(parts[1])

	if minPart == "" || !isDigits(minPart) {
		return 0, fmt.Errorf("minutes must be a number")
	}
	if secPart == "" {
		secPart = "0"
	}
	if !isDigits(secPart) {
		return 0, fmt.Errorf("seconds must be a number")
	}
	if len(secPart) > 2 {
		return 0, fmt.Errorf("seconds must be 0-59")
	}

	mins, err := strconv.Atoi(minPart)
	if err != nil || mins < 0 {
		return 0, fmt.Errorf("invalid minutes value")
	}
	secs, err := strconv.Atoi(secPart)
	if err != nil || secs < 0 || secs > 59 {
		return 0, fmt.Errorf("seconds must be 0-59")
	}

	return time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func formatJumpClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Seconds())
	mm := total / 60
	ss := total % 60
	return fmt.Sprintf("%02d:%02d", mm, ss)
}
