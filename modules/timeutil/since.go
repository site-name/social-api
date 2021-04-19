package timeutil

import (
	"fmt"
	"html/template"
	"math"
	"strings"
	"time"

	"github.com/sitename/sitename/modules/setting"

	"github.com/unknwon/i18n"
)

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func round(s float64) int64 {
	return int64(math.Round(s))
}

func computeTimeDiffFloor(diff int64, lang string) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = i18n.Tr(lang, "tool.now")
	case diff < 2:
		diff = 0
		diffStr = i18n.Tr(lang, "tool.1s")
	case diff < 1*Minute:
		diffStr = i18n.Tr(lang, "tool.seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = i18n.Tr(lang, "tool.1m")
	case diff < 1*Hour:
		diffStr = i18n.Tr(lang, "tool.minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = i18n.Tr(lang, "tool.1h")
	case diff < 1*Day:
		diffStr = i18n.Tr(lang, "tool.hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = i18n.Tr(lang, "tool.1d")
	case diff < 1*Week:
		diffStr = i18n.Tr(lang, "tool.days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = i18n.Tr(lang, "tool.1w")
	case diff < 1*Month:
		diffStr = i18n.Tr(lang, "tool.weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = i18n.Tr(lang, "tool.1mon")
	case diff < 1*Year:
		diffStr = i18n.Tr(lang, "tool.months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = i18n.Tr(lang, "tool.1y")
	default:
		diffStr = i18n.Tr(lang, "tool.years", diff/Year)
		diff -= (diff / Year) * Year
	}
	return diff, diffStr
}

func computeTimeDiff(diff int64, lang string) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = i18n.Tr(lang, "tool.now")
	case diff < 2:
		diff = 0
		diffStr = i18n.Tr(lang, "tool.1s")
	case diff < 1*Minute:
		diffStr = i18n.Tr(lang, "tool.seconds", diff)
		diff = 0

	case diff < Minute+Minute/2:
		diff -= 1 * Minute
		diffStr = i18n.Tr(lang, "tool.1m")
	case diff < 1*Hour:
		minutes := round(float64(diff) / Minute)
		if minutes > 1 {
			diffStr = i18n.Tr(lang, "tool.minutes", minutes)
		} else {
			diffStr = i18n.Tr(lang, "tool.1m")
		}
		diff -= diff / Minute * Minute

	case diff < Hour+Hour/2:
		diff -= 1 * Hour
		diffStr = i18n.Tr(lang, "tool.1h")
	case diff < 1*Day:
		hours := round(float64(diff) / Hour)
		if hours > 1 {
			diffStr = i18n.Tr(lang, "tool.hours", hours)
		} else {
			diffStr = i18n.Tr(lang, "tool.1h")
		}
		diff -= diff / Hour * Hour

	case diff < Day+Day/2:
		diff -= 1 * Day
		diffStr = i18n.Tr(lang, "tool.1d")
	case diff < 1*Week:
		days := round(float64(diff) / Day)
		if days > 1 {
			diffStr = i18n.Tr(lang, "tool.days", days)
		} else {
			diffStr = i18n.Tr(lang, "tool.1d")
		}
		diff -= diff / Day * Day

	case diff < Week+Week/2:
		diff -= 1 * Week
		diffStr = i18n.Tr(lang, "tool.1w")
	case diff < 1*Month:
		weeks := round(float64(diff) / Week)
		if weeks > 1 {
			diffStr = i18n.Tr(lang, "tool.weeks", weeks)
		} else {
			diffStr = i18n.Tr(lang, "tool.1w")
		}
		diff -= diff / Week * Week

	case diff < 1*Month+Month/2:
		diff -= 1 * Month
		diffStr = i18n.Tr(lang, "tool.1mon")
	case diff < 1*Year:
		months := round(float64(diff) / Month)
		if months > 1 {
			diffStr = i18n.Tr(lang, "tool.months", months)
		} else {
			diffStr = i18n.Tr(lang, "tool.1mon")
		}
		diff -= diff / Month * Month

	case diff < Year+Year/2:
		diff -= 1 * Year
		diffStr = i18n.Tr(lang, "tool.1y")
	default:
		years := round(float64(diff) / Year)
		if years > 1 {
			diffStr = i18n.Tr(lang, "tool.years", years)
		} else {
			diffStr = i18n.Tr(lang, "tool.1y")
		}
		diff -= (diff / Year) * Year
	}
	return diff, diffStr
}

// MinutesToFriendly returns a user friendly string with number of minutes
// converted to hours and minutes.
func MinutesToFriendly(minutes int, lang string) string {
	duration := time.Duration(minutes) * time.Minute
	return TimeSincePro(time.Now().Add(-duration), lang)
}

// TimeSincePro calculates the time interval and generate full user-friendly string.
func TimeSincePro(then time.Time, lang string) string {
	return timeSincePro(then, time.Now(), lang)
}

func timeSincePro(then, now time.Time, lang string) string {
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return i18n.Tr(lang, "tool.future")
	}
	if diff == 0 {
		return i18n.Tr(lang, "tool.now")
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiffFloor(diff, lang)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

func timeSince(then, now time.Time, lang string) string {
	return timeSinceUnix(then.Unix(), now.Unix(), lang)
}

func timeSinceUnix(then, now int64, lang string) string {
	lbl := "tool.ago"
	diff := now - then
	if then > now {
		lbl = "tool.from_now"
		diff = then - now
	}
	if diff <= 0 {
		return i18n.Tr(lang, "tool.now")
	}

	_, diffStr := computeTimeDiff(diff, lang)
	return i18n.Tr(lang, lbl, diffStr)
}

// RawTimeSince retrieves i18n key of time since t
func RawTimeSince(t time.Time, lang string) string {
	return timeSince(t, time.Now(), lang)
}

// TimeSince calculates the time interval and generate user-friendly string.
func TimeSince(then time.Time, lang string) template.HTML {
	return htmlTimeSince(then, time.Now(), lang)
}

func htmlTimeSince(then, now time.Time, lang string) template.HTML {
	return template.HTML(fmt.Sprintf(`<span class="time-since" title="%s">%s</span>`,
		then.In(setting.DefaultUILocation).Format(GetTimeFormat(lang)),
		timeSince(then, now, lang)))
}

// TimeSinceUnix calculates the time interval and generate user-friendly string.
func TimeSinceUnix(then TimeStamp, lang string) template.HTML {
	return htmlTimeSinceUnix(then, TimeStamp(time.Now().Unix()), lang)
}

func htmlTimeSinceUnix(then, now TimeStamp, lang string) template.HTML {
	return template.HTML(fmt.Sprintf(`<span class="time-since" title="%s">%s</span>`,
		then.FormatInLocation(GetTimeFormat(lang), setting.DefaultUILocation),
		timeSinceUnix(int64(then), int64(now), lang)))
}
