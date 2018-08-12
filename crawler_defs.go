package wasgeit

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	dateRe     = regexp.MustCompile(`(\d{2}.\d{2}.\d{2})`)
	dateTimeRe = regexp.MustCompile(`(\d{1,2}.\d{1,2} \d{4}) - Doors: (\d{2}:\d{2})`)
	timeRe     = regexp.MustCompile(`\d{2}:\d{2}`)
	roessliRe  = regexp.MustCompile(`\d{1,2}. \pL{3} \d{4} \d{2}:\d{2}`)
)

var kairoConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: "article[id]",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".concerts_date").Parent().Text()
		timeString := timeRe.FindString(rawDateTimeString)
		return rawDateTimeString[3:13] + timeString
	},
	TitleSelector: "h1",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return fmt.Sprint(venue.URL, "#", eventSelection.AttrOr("id", ""))
	}}

var dachstockConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".event.event-list",
	TimeFormat:    "2.1 200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".event-date").Text()
		captures := dateTimeRe.FindStringSubmatch(rawDateTimeString)
		return captures[1] + captures[2]
	},
	TitleSelector: "h3",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.AttrOr("data-url", venue.URL)
	}}

var turnhalleConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".event",
	TimeFormat:    "02. 01. 0615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("h4").Text()
		dateString := rawDateTimeString[4:14]
		matches := timeRe.FindAllStringSubmatch(rawDateTimeString, 2)
		var timeString string
		if len(matches) > 0 && len(matches[0]) == 1 {
			timeString = matches[0][0]
		}
		return dateString + timeString
	},
	TitleSelector: "h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return fmt.Sprint(venue.URL, eventSelection.Find("a").AttrOr("href", ""))
	}}

var brasserieLorraineConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".type-tribe_events",
	TimeFormat:    "January 2",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".tribe-event-schedule-details").Text()
		tokens := strings.Split(rawDateTimeString, " @ ")
		return strings.TrimSpace(tokens[0])
	},
	TitleSelector: ".tribe-events-list-event-title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("h2 > a").AttrOr("href", venue.URL)
	}}

var kofmehlConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".events__element",
	TimeFormat:    "02.01",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("time").Text()
		dateString := rawDateTimeString[3:8]
		return dateString
	},
	TitleSelector: ".events__title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("a.events__link").AttrOr("href", venue.URL)
	}}

var kiffConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".programm-grid a:not(.teaserlink)",
	TimeFormat:    "2 Jan",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event-date").Text()[3:]
	},
	TitleSelector: ".event-title-wrapper > h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Attr("href"); exists {
			base, _ := url.Parse(venue.URL)
			relative, _ := url.Parse(href)
			return base.ResolveReference(relative).String()
		}
		return venue.URL
	}}

var coqDorConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: "#main table:not(.shows)",
	TimeFormat:    "02.01.0615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("td.list_first").Text()
		dateString := dateRe.FindString(rawDateTimeString)
		rawTimeString := eventSelection.Find("div.entry").Text()
		timeString := timeRe.FindString(rawTimeString)
		return dateString + timeString
	},
	TitleSelector: "td.list_second h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("td.list_second h2 a").AttrOr("href", venue.URL)
	}}

var iscConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".page_programm a.event_preview",
	TimeFormat:    "02.01.",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event_title_date").Text()
	},
	TitleSelector: ".event_title_title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.AttrOr("href", venue.URL)
	}}

var mahoganyHallConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".view-konzerte .views-row",
	TimeFormat:    "02. January 2006|15.04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dateTimeString := eventSelection.Find(".concert-tueroeffnung").Text()
		dateTimeString = StripSomeWhiteSpaces(dateTimeString)
		dateTimeString = strings.Split(dateTimeString, ", ")[1]
		dateTimeString = strings.Split(dateTimeString, "Uhr")[0]
		if dateTimeString[1] == '.' {
			return "0" + dateTimeString
		}
		return dateTimeString
	},
	TitleSelector: ".views-field-title h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".views-field-title h2 a").Attr("href"); exists {
			base, _ := url.Parse(venue.URL)
			relative, _ := url.Parse(href)
			return base.ResolveReference(relative).String()
		}
		return venue.URL
	}}

var heitereFahneConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".events .event",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".date + .time").Parent().Text()
		rawDateTimeString = StripSomeWhiteSpaces(rawDateTimeString)
		rawDateTimeString = strings.TrimSpace(rawDateTimeString)
		return rawDateTimeString[3:13] + rawDateTimeString[33:38]
	},
	TitleSelector: ".alpha.omega.text .inner h2 a",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		if href, ok := eventSelection.Find(".alpha.omega.text .inner h2 a").Attr("href"); ok {
			base, _ := url.Parse(venue.URL)
			relative, _ := url.Parse(href)
			return base.ResolveReference(relative).String()
		}
		return venue.URL
	}}

var onoConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".EventItem",
	TimeFormat:    "02.01.0615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".EventInfo.subnav").Text()
		rawDateTimeString = wrp.Replace(rawDateTimeString)
		dateString := rawDateTimeString[3:11]
		timeString := timeRe.FindString(rawDateTimeString)
		return dateString + timeString
	},
	TitleSelector: ".EventTextTitle",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		href := eventSelection.Find(".EventImage a").AttrOr("href", "")
		base, _ := url.Parse(venue.URL)
		relative, _ := url.Parse(href)
		return base.ResolveReference(relative).String()
	}}

var martaConfig = HTMLConfig{
	IsSameEvent:   hasSameTitleAndDate,
	EventSelector: "table.music tbody tr",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dateString := eventSelection.Find("td:nth-child(1)").Text()
		rawTimeString := eventSelection.Find("td:nth-child(4)").Text()
		timeString := timeRe.FindString(rawTimeString)
		return dateString + timeString
	},
	TitleSelector: "td:nth-child(3) p",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		href := eventSelection.Find(".EventImage a").AttrOr("href", "")
		return fmt.Sprint(venue.URL, href)
	}}

var bierhuebeliConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: "ul.bh-event-list.all-events li",
	TimeFormat:    "02.01.06",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawTimeString := eventSelection.Find(".evendates").Text()
		return rawTimeString[8:16]
	},
	TitleSelector: ".eventlink a",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find(".eventlink a").AttrOr("href", venue.URL)
	}}

var dampfzentraleConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: "article .agenda-container",
	TimeFormat:    "2.1.15:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		article := eventSelection.Parent().Parent()
		month, _ := article.Attr("data-month")
		day, _ := article.Attr("data-date")
		dateString := fmt.Sprintf("%s.%s.", day, month)
		timeString := strings.TrimSpace(eventSelection.Find(".agenda-details .span1").Text())
		return dateString + timeString
	},
	TitleSelector: "h1.agenda-title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		id := eventSelection.Parent().AttrOr("id", "")
		return fmt.Sprintf("%s#%s", venue.URL, id)
	}}

var roessliConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".rossli-events .event",
	TimeFormat:    "2. Jan 2006 15:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dt := eventSelection.Find("time.event-date").AttrOr("datetime", "")
		replaced := strings.Replace(dt, "Mrz", "Mär", -1)
		return roessliRe.FindString(replaced)
		// return dt[4:21]
	},
	TitleSelector: "h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("a").AttrOr("href", venue.URL)
	}}

var souslepontConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".sous-le-pont-programm .event",
	TimeFormat:    "2. Jan 2006 15:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dt := eventSelection.Find("time.event-date").AttrOr("datetime", "")
		replaced := strings.Replace(dt, "Mrz", "Mär", -1)
		return roessliRe.FindString(replaced)
	},
	TitleSelector: "h2",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("a").AttrOr("href", venue.URL)
	}}

var lesAmisConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".cff-event",
	TimeFormat:    "Jan 2, 3:04pm",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dateStr := eventSelection.Find(".cff-date > .cff-start-date").Text()
		return strings.Replace(dateStr, "Mrz", "Mär", -1)
	},
	TitleSelector: ".cff-event-title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		id := eventSelection.AttrOr("id", "")
		return fmt.Sprintf("%s#%s", venue.URL, id)
	}}

var mokkaConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".event-month > a",
	TimeFormat:    "Mon. 02. Jan.",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find("div.date").Text()
	},
	TitleSelector: "div.title-section",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.AttrOr("href", venue.URL)
	}}

var muehleHunzikenConfig = HTMLConfig{
	IsSameEvent:   hasSameUrl,
	EventSelector: ".event-list-item",
	TimeFormat:    "02.01.2006",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dateStr := eventSelection.Find(".event-date").Text()
		return dateStr[4:16]
	},
	TitleSelector: ".event-title",
	LinkBuilder: func(venue Venue, eventSelection *goquery.Selection) string {
		return eventSelection.Find("a").AttrOr("href", venue.URL)
	}}

// TODO
// http://wartsaal-kaffee.ch/veranstaltungen/
// http://www.cafete.ch/
// http://www.schlachthaus.ch/spielplan/index.php
// https://www.effinger.ch/events/
// http://dynamo.ch/veranstaltungen?field_event_type_tid=1&field_event_zeitraum_value_1[value][month]=1&field_event_zeitraum_value_1[value][year]=2018

func RegisterAllHTMLCrawlers(st *Store) {
	registerHTMLCrawler("kairo", kairoConfig, st)
	registerHTMLCrawler("dachstock", dachstockConfig, st)
	registerHTMLCrawler("turnhalle", turnhalleConfig, st)
	registerHTMLCrawler("brasserie-lorraine", brasserieLorraineConfig, st)
	registerHTMLCrawler("kofmehl", kofmehlConfig, st)
	registerHTMLCrawler("kiff", kiffConfig, st)
	registerHTMLCrawler("coq-d-or", coqDorConfig, st)
	registerHTMLCrawler("isc", iscConfig, st)
	registerHTMLCrawler("mahogany-hall", mahoganyHallConfig, st)
	registerHTMLCrawler("heitere-fahne", heitereFahneConfig, st)
	registerHTMLCrawler("ono", onoConfig, st)
	registerHTMLCrawler("marta", martaConfig, st)
	registerHTMLCrawler("bierhuebeli", bierhuebeliConfig, st)
	registerHTMLCrawler("dampfzentrale", dampfzentraleConfig, st)
	registerHTMLCrawler("roessli", roessliConfig, st)
	registerHTMLCrawler("sous-le-pont", souslepontConfig, st)
	registerHTMLCrawler("les-amis", lesAmisConfig, st)
	registerHTMLCrawler("mokka", mokkaConfig, st)
	registerHTMLCrawler("muehle-hunziken", muehleHunzikenConfig, st)
}

func registerHTMLCrawler(shortName string, config HTMLConfig, st *Store) {
	RegisterCrawler(shortName, &HTMLCrawler{config: config, venue: st.GetVenue(shortName)})
}
