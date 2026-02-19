package models

type Unit string

const (
	UnitSecond Unit = "second"
	UnitMinute Unit = "minute"
	UnitHour   Unit = "hour"
	UnitDay    Unit = "day"
	UnitWeek   Unit = "week"
	UnitMonth  Unit = "month"
)

type Request struct {
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}
type Job struct {
	User_name         string  `json:"user_name"`
	Name              string  `json:"name"`
	Target            string  `json:"target"`
	Request           Request `json:"request"`
	Expected_response string  `json:"expected_response"`
	Schedule          string  `json:"schedule"`
}

type Schedule struct {
	Every   int    `json:"every"`
	Unit    Unit   `json:"unit"`
	At      string `json:"at,omitempty"`
	Weekday int    `json:"weekday,omitempty"`
	Day     int    `json:"day,omitempty"`
}
