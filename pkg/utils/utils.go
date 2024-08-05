package utils

import (
	"regexp"
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
)

func HumanTimeToUnix(humanTime string) int64 {
	if humanTime == "never" {
		return -1
	}
	duration := config.ParseDuration(humanTime)
	if config.ParseDuration(config.Vars.MaxExpiryTime) < duration {
		return time.Now().Add(time.Duration(config.OneWeek)).Unix()
	}
	return time.Now().Add(duration).Unix()
}

func IsContentJustUrl(content string) int {
	if regexp.MustCompile(`^(?:http|https|magnet):\/\/[^\s/$.?#].[^\s]*$`).MatchString(content) {
		return 1
	}
	return 0
}

func Ternary(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
