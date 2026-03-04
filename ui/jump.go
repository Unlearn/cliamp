package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseJumpTarget parses an absolute jump target from track start
func parseJumpTarget(raw string) (time.Duration, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("enter time like 10 or 58:05")
	}

	minPart, secPart, hasSep := strings.Cut(s, ":")
	if !hasSep {
		secs, err := parseWholeSeconds(minPart)
		if err != nil {
			return 0, err
		}
		return time.Duration(secs) * time.Second, nil
	}
	if strings.Contains(secPart, ":") {
		return 0, fmt.Errorf("use mm:ss format")
	}

	minPart = normalizeClockPart(minPart)
	secPart = normalizeClockPart(secPart)

	mins, err := parseMinutesPart(minPart)
	if err != nil {
		return 0, err
	}
	secs, err := parseSecondsPart(secPart)
	if err != nil {
		return 0, err
	}

	return time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second, nil
}

func normalizeClockPart(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "0"
	}
	return s
}

func parseWholeSeconds(s string) (int, error) {
	secs, err := parsePositiveInt(strings.TrimSpace(s), "seconds")
	if err != nil {
		return 0, err
	}
	return secs, nil
}

func parseMinutesPart(s string) (int, error) {
	mins, err := parsePositiveInt(s, "minutes")
	if err != nil {
		return 0, err
	}
	return mins, nil
}

func parseSecondsPart(s string) (int, error) {
	if len(s) > 2 {
		return 0, fmt.Errorf("seconds must be 0-59")
	}

	secs, err := parsePositiveInt(s, "seconds")
	if err != nil {
		return 0, err
	}
	if secs > 59 {
		return 0, fmt.Errorf("seconds must be 0-59")
	}
	return secs, nil
}

func parsePositiveInt(s, label string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("%s must be a number", label)
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("%s must be a number", label)
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", label)
	}
	return v, nil
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
