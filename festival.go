package wasgeit

import "time"

type Festival struct {
	Id           int64         `json:"id"`
	Url          string        `json:"url"`
	Title        string        `json:"title"`
	Location     string        `json:"location"`
	DateStart    time.Time     `json:"date_start"`
	DateEnd      time.Time     `json:"date_end"`
	OpeningTimes []OpeningTime `json:"opening_times"`
}

type OpeningTime struct {
	Days  string `json:"days"`
	Start string `json:"start"`
	End   string `json:"end"`
}
