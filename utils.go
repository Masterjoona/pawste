package main

import (
	"fmt"
	"time"
)

func SubmitToPaste(submit Submit, pasteName string, hashedPassword string) Paste {
	return Paste{
		PasteName:      pasteName,
		Expire:         HumanTimeToSQLTime(submit.Expiration),
		Privacy:        submit.Privacy,
		BurnAfter:      submit.BurnAfter,
		Content:        submit.Text,
		Syntax:         submit.Syntax,
		HashedPassword: hashedPassword,
	}
}

func HumanTimeToSQLTime(humanTime string) string {
	var duration time.Duration
	switch humanTime {
	case "10min":
		duration = 10 * time.Minute
	case "1min":
		duration = 1 * time.Minute
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	case "72h":
		duration = 72 * time.Hour
	case "1w":
		duration = 7 * 24 * time.Hour
	case "never":
		duration = 100 * 365 * 25 // cope if you're still using this in 100 years
	default:
		duration = 7 * 24 * time.Hour
	}
	return fmt.Sprintf("%d", time.Now().Add(duration).Unix())
}
