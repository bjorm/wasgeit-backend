package wasgeit

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
)

var (
	location, _ = time.LoadLocation("Europe/Zurich")
)

type HTMLConfig struct {
	EventSelector     string
	TitleSelector     string
	GetDateTimeString func(*goquery.Selection) string
	TimeFormat        string
	LinkBuilder       func(Venue, *goquery.Selection) string
	IsSameEvent       func(ev1, ev2 Event) bool
}

type HTMLCrawler struct {
	venue  Venue
	dom    *goquery.Document
	config HTMLConfig
}

func (cr *HTMLCrawler) Name() string {
	return cr.venue.ShortName
}

func (cr *HTMLCrawler) URL() string {
	return cr.venue.URL
}

func (cr *HTMLCrawler) IsSame(ev1, ev2 Event) bool {
	return cr.config.IsSameEvent(ev1, ev2)
}

func (cr *HTMLCrawler) Fetch() error {
	dom, err := goquery.NewDocument(cr.venue.URL)
	if err != nil {
		return err
	}
	cr.dom = dom
	return nil
}

func (cr *HTMLCrawler) LoadFrom(r io.Reader) error {
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	cr.dom = dom
	return nil
}

func (cr *HTMLCrawler) GetEvents() ([]Event, []error) {
	var evs []Event
	var errors []error

	cr.dom.Find(cr.config.EventSelector).Each(func(_ int, eventSelection *goquery.Selection) {
		re := HTMLEvent{s: eventSelection, c: cr.config, v: cr.venue}
		datetime, err := re.dateTime()
		if err != nil {
			errors = append(errors, err)
		} else if datetime.After(time.Now()) {
			evs = append(evs, Event{DateTime: datetime, Title: re.title(), URL: re.url(), Venue: cr.venue})
		}
	})

	return evs, errors
}

type HTMLEvent struct {
	s *goquery.Selection
	c HTMLConfig
	v Venue
}

func (e *HTMLEvent) title() string {
	tr := strings.TrimSpace(e.s.Find(e.c.TitleSelector).Text())
	return StripLineBreaks(tr)
}

func (e *HTMLEvent) url() string {
	return e.c.LinkBuilder(e.v, e.s)
}

func (e *HTMLEvent) dateTime() (time.Time, error) {
	timeStr := e.c.GetDateTimeString(e.s)

	if timeStr == "" {
		return time.Time{}, fmt.Errorf("Time selector yielded empty string")
	}

	timeStr = strings.TrimSpace(timeStr)
	eventTime, timeParseError := monday.ParseInLocation(e.c.TimeFormat, timeStr, location, monday.LocaleDeDE)

	if timeParseError != nil {
		return time.Time{}, timeParseError
	}

	// TODO maybe move this into post-process method
	// Some sites publish their events without specifying a year, we assume they take place this year.
	if eventTime.Year() == 0 {
		eventTime = eventTime.AddDate(time.Now().Year(), 0, 0)
	}

	return eventTime, nil
}

func returnStringSlice(start int, end int) func(string) string {
	return func(toSlice string) string {
		return toSlice[start:end]
	}
}

func StripLineBreaks(s string) string {
	tokens := strings.Split(s, "\n")
	if len(tokens) == 1 {
		return s
	}

	var n []string
	for _, t := range tokens {
		n = append(n, strings.TrimSpace(t))
	}

	return strings.Join(n, " ")
}

var wrp = strings.NewReplacer("\u2009", "", "\u00a0", "", "\n", "", "\t", "")

// TODO find out why there are still whitespaces in some titles
// StripSomeWhiteSpaces strips the following whitespaces: \u00a0, \n, \t
func StripSomeWhiteSpaces(toStrip string) string {
	return wrp.Replace(toStrip)
}

func hasSameUrl(ev1, ev2 Event) bool {
	return ev1.URL == ev2.URL
}

func hasSameTitleAndDate(ev1, ev2 Event) bool {
	return ev1.Title == ev2.Title && ev1.DateTime.Equal(ev2.DateTime)
}
