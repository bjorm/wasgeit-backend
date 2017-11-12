package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
)

var (
	location, _ = time.LoadLocation("Europe/Zurich")
)

// Event describes an event taking place in a Venue
type Event struct {
	ID       int64
	Title    string
	DateTime time.Time
	URL      string
	Venue    Venue
}

// Venue describes a place where Events take place
type Venue struct {
	ID   int64
	Name string
	URL  string
}

// Crawler describes  a crawler
type Crawler interface {
	Crawl() ([]Event, error)
	Venue() Venue
}

type HTMLCrawler struct {
	venue               Venue
	eventSelector       string
	titleSelector       string
	timeSelector        string
	timeParsePreProcess func(string) string
	timeFormat          string
	linkBuilder         func(*HTMLCrawler, *goquery.Selection) string
}

func (cr *HTMLCrawler) Venue() Venue {
	return cr.venue
}

func (cr *HTMLCrawler) Crawl() (events []Event, err error) {
	document, loadError := goquery.NewDocument(cr.venue.URL)
	if loadError != nil {
		return events, loadError
	}

	document.Find(cr.eventSelector).Each(func(_ int, eventSelection *goquery.Selection) {
		title := getTrimmedText(eventSelection, cr.titleSelector)
		time, err := cr.getEventTime(eventSelection)

		if err == nil {
			linkURL := cr.linkBuilder(cr, eventSelection)
			event := Event{DateTime: *time, Title: title, URL: linkURL, Venue: cr.venue}
			events = append(events, event)
		} else {
			glog.Warningf("Skipped %q because of: %q", title, err)
		}
	})

	return events, nil
}

func (cr *HTMLCrawler) getEventTime(event *goquery.Selection) (*time.Time, error) {
	timeStr := getTrimmedText(event, cr.timeSelector)

	if timeStr == "" {
		return nil, fmt.Errorf("Time selector yielded empty string")
	}

	timeStr = cr.timeParsePreProcess(timeStr)
	timeStr = strings.TrimSpace(timeStr)
	eventTime, timeParseError := monday.ParseInLocation(cr.timeFormat, timeStr, location, monday.LocaleDeDE)

	if timeParseError != nil {
		return nil, timeParseError
	}

	// TODO maybe move this into post-process method
	// Some sites publish their events without specifying a year, we assume they take place this year.
	if eventTime.Year() == 0 {
		eventTime = eventTime.AddDate(time.Now().Year(), 0, 0)
	}

	return &eventTime, nil
}

func getTrimmedText(selection *goquery.Selection, selector string) string {
	return strings.TrimSpace(selection.Find(selector).Text())
}

func returnStringSlice(start int, end int) func(string) string {
	return func(toSlice string) string {
		return toSlice[start:end]
	}
}
