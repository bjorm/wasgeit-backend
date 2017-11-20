package wasgeit

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
	ID        int64
	ShortName string
	Name      string
	URL       string
}

// Crawler describes  a crawler
type Crawler interface {
	Crawl() ([]Event, error)
	Venue() Venue
}

type HTMLCrawler struct {
	venue             Venue
	EventSelector     string
	TitleSelector     string
	GetDateTimeString func(*goquery.Selection) string
	TimeFormat        string
	LinkBuilder       func(*HTMLCrawler, *goquery.Selection) string
}

func (cr *HTMLCrawler) Venue() Venue {
	return cr.venue
}

func (cr *HTMLCrawler) Crawl() (events []Event, err error) {
	document, loadError := goquery.NewDocument(cr.venue.URL)
	if loadError != nil {
		return events, loadError
	}

	document.Find(cr.EventSelector).Each(func(_ int, eventSelection *goquery.Selection) {
		title := getTrimmedText(eventSelection, cr.TitleSelector)
		time, err := cr.GetEventTime(eventSelection)
		if err == nil {
			linkURL := cr.LinkBuilder(cr, eventSelection)
			event := Event{DateTime: *time, Title: title, URL: linkURL, Venue: cr.venue}
			events = append(events, event)
		} else {
			glog.Warningf("Skipped %q because of: %q", title, err)
		}
	})

	return events, nil
}

func (cr *HTMLCrawler) GetEventTime(event *goquery.Selection) (*time.Time, error) {
	timeStr := cr.GetDateTimeString(event)

	if timeStr == "" {
		return nil, fmt.Errorf("Time selector yielded empty string")
	}

	timeStr = strings.TrimSpace(timeStr)
	eventTime, timeParseError := monday.ParseInLocation(cr.TimeFormat, timeStr, location, monday.LocaleDeDE)

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

func returnStringSlice(start int, end int) func(string) string {
	return func(toSlice string) string {
		return toSlice[start:end]
	}
}

func getTrimmedText(selection *goquery.Selection, selector string) string {
	return strings.TrimSpace(selection.Find(selector).Text())
}

var wrp = strings.NewReplacer("\u2009", "", "\u00a0", "", "\n", "", "\t", "")

// StripSomeWhiteSpaces strips the following whitespaces: \u00a0, \n, \t
func StripSomeWhiteSpaces(toStrip string) string {
	return wrp.Replace(toStrip)
}
