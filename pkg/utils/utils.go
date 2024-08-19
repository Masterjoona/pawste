package utils

import (
	"regexp"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

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

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

var stripRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
var trans = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func NormalizeFilename(str string) string {
	//return stripRegex.ReplaceAllString(str, "")
	t, _, _ := transform.String(trans, str)
	t = stripRegex.ReplaceAllString(t, "")
	return t
}
