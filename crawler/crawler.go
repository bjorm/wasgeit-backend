package crawler

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
}

type HTMLCrawler struct {
	venue         Venue
	eventSelector string
	titleSelector string
	timeSelector  string
	timeFormat    string
	linkBuilder   func(*HTMLCrawler, *goquery.Selection) string
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

	// TODO improve readability
	timeStr = Stencil(cr.timeFormat, timeStr)
	timeFormat := stripDashes(cr.timeFormat)

	eventTime, timeParseError := monday.ParseInLocation(timeFormat, timeStr, location, monday.LocaleDeDE)

	if timeParseError != nil {
		return nil, timeParseError
	}

	// Some sites publish their events without specifying a year, we assume they take place this year.
	if eventTime.Year() == 0 {
		eventTime = eventTime.AddDate(time.Now().Year(), 0, 0)
	}

	return &eventTime, nil
}

func getTrimmedText(selection *goquery.Selection, selector string) string {
	return strings.TrimSpace(selection.Find(selector).Text())
}

var kairoCrawler = HTMLCrawler{
	venue:         Venue{ID: 1, Name: "Cafe Kairo", URL: "http://www.cafe-kairo.ch/kultur"},
	eventSelector: "article",
	timeSelector:  "time",
	timeFormat:    "---02.01.2006",
	titleSelector: "h1",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if id, exists := eventSelection.Attr("id"); exists {
			return fmt.Sprintf("%s#%s", crawler.venue.URL, id)
		}
		return crawler.venue.URL
	}}

var dachstockCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Dachstock", URL: "http://www.dachstock.ch"},
	eventSelector: ".em-eventlist-event",
	timeSelector:  ".em-eventlist-date",
	timeFormat:    "-----2.1 2006  ", // TODO fix hack with trailing space
	titleSelector: "h3",
	linkBuilder: func(crawler *HTMLCrawler, _ *goquery.Selection) string {
		return crawler.venue.URL
	}}

var turnhalleCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Turnhalle", URL: "http://www.turnhalle.ch"},
	eventSelector: ".event",
	timeSelector:  "h4",
	timeFormat:    "----02. 01. 06",
	titleSelector: "h2",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("a").Attr("href"); exists {
			return fmt.Sprintf("%s%s", crawler.venue.URL, href)
		}
		return crawler.venue.URL
	}}

var brasserieLorraineCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Brasserie Lorraine", URL: "http://brasserie-lorraine.ch/?post_type=tribe_events"},
	eventSelector: ".type-tribe_events",
	timeSelector:  ".tribe-event-date-start",
	timeFormat:    "January 02 ", // TODO fix hack with trailing space
	titleSelector: ".tribe-events-list-event-title",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("h2 > a").Attr("href"); exists {
			return fmt.Sprintf("%s%s", crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var Crawlers = []Crawler{
	&kairoCrawler,
	&dachstockCrawler,
	&turnhalleCrawler,
	&brasserieLorraineCrawler}
