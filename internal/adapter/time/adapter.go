package time

import (
	"auth-service/internal/pkg/consts"
	"time"
)

type (
	Adapter interface {
		Now() time.Time
		AddTime(t time.Time, d time.Duration) time.Time
		MillisecondsToTime(milliseconds int64) time.Time
		Locale() *time.Location
		LocaleOffsetMilli() int64
	}

	adapter struct {
		localeOffsetMilli int64
		locale            *time.Location
	}
)

func NewAdapter(locale int64) Adapter {
	return &adapter{
		localeOffsetMilli: locale * consts.HourMilli,
		locale:            time.FixedZone("MSC", int(locale)*3600),
	}
}

func (a *adapter) Now() time.Time {
	return time.Now().In(a.locale)
}

func (a *adapter) AddTime(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

func (a *adapter) Locale() *time.Location {
	return a.locale
}

func (a *adapter) LocaleOffsetMilli() int64 {
	return a.localeOffsetMilli
}

func (a *adapter) MillisecondsToTime(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds).In(a.locale)
}
