package wasgeit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var dateTimeRe = regexp.MustCompile(`(\d{1,2}.\d{1,2} \d{4}) - Doors: (\d{2}:\d{2})`)
var timeRe = regexp.MustCompile(`\d{2}:\d{2}`)

var kairoCrawler = HTMLCrawler{
	venue:         Venue{ID: 1, Name: "Cafe Kairo", ShortName: "kairo", URL: "http://www.cafe-kairo.ch/kultur"},
	EventSelector: "article[id]",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".concerts_date").Parent().Text()
		return rawDateTimeString[3:13] + rawDateTimeString[19:24]
	},
	TitleSelector: "h1",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if id, exists := eventSelection.Attr("id"); exists {
			return fmt.Sprintf("%s#%s", crawler.venue.URL, id)
		}
		return crawler.venue.URL
	}}

var dachstockCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Dachstock", ShortName: "dachstock", URL: "http://www.dachstock.ch"},
	EventSelector: ".em-eventlist-event",
	TimeFormat:    "2.1 200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".em-eventlist-date").Text()
		captures := dateTimeRe.FindStringSubmatch(rawDateTimeString)
		return captures[1] + captures[2]
	},
	TitleSelector: "h3",
	LinkBuilder: func(crawler *HTMLCrawler, _ *goquery.Selection) string {
		return crawler.venue.URL
	}}

var turnhalleCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Turnhalle", ShortName: "turnhalle", URL: "http://www.turnhalle.ch"},
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
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("a").Attr("href"); exists {
			return fmt.Sprintf("%s%s", crawler.venue.URL, href)
		}
		return crawler.venue.URL
	}}

var brasserieLorraineCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Brasserie Lorraine", ShortName: "brasserie-lorraine", URL: "http://brasserie-lorraine.ch/?post_type=tribe_events"},
	EventSelector: ".type-tribe_events",
	TimeFormat:    "January 2",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".tribe-event-date-start").Text()
		dateString := rawDateTimeString[0:11]
		return strings.TrimSpace(dateString)
	},
	TitleSelector: ".tribe-events-list-event-title",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("h2 > a").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var kofmehlCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Kofmehl", ShortName: "kofmehl", URL: "http://www.kofmehl.net"},
	EventSelector: ".events__element",
	TimeFormat:    "02.01",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("time").Text()
		dateString := rawDateTimeString[3:8]
		return dateString
	},
	TitleSelector: ".events__title",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("a.events__link").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var kiffCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Kiff", ShortName: "kiff", URL: "http://www.kiff.ch"},
	EventSelector: ".programm-grid a:not(.teaserlink)",
	TimeFormat:    "2 Jan",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event-date").Text()[3:]
	},
	TitleSelector: ".event-title-wrapper > h2",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Attr("href"); exists {
			return fmt.Sprintf("%s%s", crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var coqDorCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Coq d'Or", ShortName: "coq-d-or", URL: "http://www.coq-d-or.ch/"},
	EventSelector: "#main table:not(.shows)",
	TimeFormat:    "02.01.0615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("td.list_first a").Text()
		dateString := strings.Split(rawDateTimeString, ", ")[1]
		rawTimeString := eventSelection.Find("div.entry").Text()
		timeString := timeRe.FindString(rawTimeString)
		return dateString + timeString
	},
	TitleSelector: "td.list_second h2",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("td.list_second h2 a").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var iscCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "ISC", ShortName: "isc", URL: "http://www.isc-club.ch/"},
	EventSelector: ".page_programm a.event_preview",
	TimeFormat:    "02.01.",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event_title_date").Text()
	},
	TitleSelector: ".event_title_title",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var mahoganyHallCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Mahogany Hall", ShortName: "mahogany-hall", URL: "https://www.mahogany.ch/konzerte"},
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
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".views-field-title h2 a").Attr("href"); exists {
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var heitereFahneCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Heitere Fahne", ShortName: "heitere-fahne", URL: "http://www.dieheiterefahne.ch/de/hauptnavigation/start/programm-31.html"},
	EventSelector: ".events .event",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".date + .time").Parent().Text()
		rawDateTimeString = StripSomeWhiteSpaces(rawDateTimeString)
		rawDateTimeString = strings.TrimSpace(rawDateTimeString)
		return rawDateTimeString[3:13] + rawDateTimeString[33:38]
	},
	TitleSelector: ".alpha.omega.text .inner h2 a",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".alpha.omega.text .inner h2 a").Attr("href"); exists {
			// TODO FIXME the extracted href starts with /. Simply concatenating with the venue URl will not work.
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var onoCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "ONO", ShortName: "ono", URL: "http://www.onobern.ch/programm-bersicht"},
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
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".EventImage a").Attr("href"); exists {
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var martaCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Cafe Marta", ShortName: "marta", URL: "http://www.cafemarta.ch/musik"},
	EventSelector: "table.music tbody tr",
	TimeFormat:    "02.01.200615:04",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		dateString := eventSelection.Find("td:nth-child(1)").Text()
		rawTimeString := eventSelection.Find("td:nth-child(4)").Text()
		timeString := timeRe.FindString(rawTimeString)
		return dateString + timeString
	},
	TitleSelector: "td:nth-child(3) p",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".EventImage a").Attr("href"); exists {
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var bierhuebeliCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Bierhuebeli", ShortName: "bierhuebeli", URL: "http://www.bierhuebeli.ch/veranstaltungen/"},
	EventSelector: "ul.bh-event-list.all-events li",
	TimeFormat:    "02.01.06",
	GetDateTimeString: func(eventSelection *goquery.Selection) string {
		rawTimeString := eventSelection.Find(".evendates").Text()
		return rawTimeString[8:16]
	},
	TitleSelector: ".eventlink a",
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".eventlink a").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}
var dampfzentraleCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Dampfzentrale", ShortName: "dampfzentrale", URL: "http://dampfzentrale.ch/programm/"},
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
	LinkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if id, exists := eventSelection.Parent().Attr("id"); exists {
			return fmt.Sprintf("%s#%s", crawler.venue.URL, id)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var HTMLCrawlers = map[string]HTMLCrawler{
	iscCrawler.venue.ShortName:               iscCrawler,
	kiffCrawler.venue.ShortName:              kiffCrawler,
	kofmehlCrawler.venue.ShortName:           kofmehlCrawler,
	kairoCrawler.venue.ShortName:             kairoCrawler,
	dachstockCrawler.venue.ShortName:         dachstockCrawler,
	coqDorCrawler.venue.ShortName:            coqDorCrawler,
	turnhalleCrawler.venue.ShortName:         turnhalleCrawler,
	brasserieLorraineCrawler.venue.ShortName: brasserieLorraineCrawler,
	mahoganyHallCrawler.venue.ShortName:      mahoganyHallCrawler,
	heitereFahneCrawler.venue.ShortName:      heitereFahneCrawler,
	onoCrawler.venue.ShortName:               onoCrawler,
	martaCrawler.venue.ShortName:             martaCrawler,
	bierhuebeliCrawler.venue.ShortName:       bierhuebeliCrawler,
	dampfzentraleCrawler.venue.ShortName:     dampfzentraleCrawler}

var Crawlers = []Crawler{
	&iscCrawler,
	&kiffCrawler,
	&kofmehlCrawler,
	&kairoCrawler,
	&coqDorCrawler,
	&dachstockCrawler,
	&turnhalleCrawler,
	&brasserieLorraineCrawler,
	&mahoganyHallCrawler,
	&heitereFahneCrawler,
	&onoCrawler,
	&martaCrawler,
	&bierhuebeliCrawler,
	&dampfzentraleCrawler}

// http://wartsaal-kaffee.ch/veranstaltungen/
// https://www.facebook.com/pg/CaffeBarSattler/events/t
// http://www.cafete.ch/
// http://www.schlachthaus.ch/spielplan/index.php
// https://www.effinger.ch/events/
// https://www.facebook.com/pg/loescherbern/events/?ref=page_internal
// https://www.facebook.com/peterflamingobern/
// roessli, sous-le-pont,
