package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var kairoCrawler = HTMLCrawler{
	venue:         Venue{ID: 1, Name: "Cafe Kairo", URL: "http://www.cafe-kairo.ch/kultur"},
	eventSelector: "article[id]",
	timeFormat:    "02.01.2006",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".concerts_date").Parent().Text()
		fmt.Printf("extracted time: %q\n", rawDateTimeString)
		return rawDateTimeString[3:13]
	},
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
	timeFormat:    "2.1 200615:04",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".em-eventlist-date").Text()
		return rawDateTimeString[5:15] + rawDateTimeString[25:29]
	},
	titleSelector: "h3",
	linkBuilder: func(crawler *HTMLCrawler, _ *goquery.Selection) string {
		return crawler.venue.URL
	}}

var turnhalleCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Turnhalle", URL: "http://www.turnhalle.ch"},
	eventSelector: ".event",
	timeFormat:    "02. 01. 06",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("h4").Text()
		dateString := rawDateTimeString[5:15]
		return dateString
	},
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
	timeFormat:    "January 02",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".tribe-event-date-start").Text()
		dateString := rawDateTimeString[0:11]
		return dateString
	},
	titleSelector: ".tribe-events-list-event-title",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("h2 > a").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var kofmehlCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Kofmehl", URL: "http://www.kofmehl.net"},
	eventSelector: ".events__element",
	timeFormat:    "02.01",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("time").Text()
		dateString := rawDateTimeString[3:8]
		return dateString
	},
	titleSelector: ".events__title",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("a.events__link").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var kiffCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Kiff", URL: "http://www.kiff.ch"},
	eventSelector: ".programm-grid a",
	timeFormat:    "2 Jan",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event-date").Text()[3:]
	},
	titleSelector: ".event-title-wrapper > h2",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Attr("href"); exists {
			return fmt.Sprintf("%s%s", crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var coqDorCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Coq d'Or", URL: "http://www.coq-d-or.ch/"},
	eventSelector: "#main table",
	timeFormat:    "02.01.06",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find("td.list_first a").Text()
		dateString := strings.Split(rawDateTimeString, ", ")[1]
		return dateString
	},
	titleSelector: "td.list_second h2",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find("td.list_second h2 a").Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var iscCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "ISC", URL: "http://www.isc-club.ch/"},
	eventSelector: ".page_programm a.event_preview",
	timeFormat:    "02.01.",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		return eventSelection.Find(".event_title_date").Text()
	},
	titleSelector: ".event_title_title",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Attr("href"); exists {
			return href
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var mahoganyHallCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Mahogany Hall", URL: "https://www.mahogany.ch/konzerte"},
	eventSelector: ".view-konzerte .views-row",
	timeSelector:  ".concert-tueroeffnung",
	timeFormat:    "02. January 2006|15.04",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		dateTimeString := eventSelection.Find("time").Text()
		dateTimeString = StripSomeWhiteSpaces(dateTimeString)
		dateTimeString = strings.Split(dateTimeString, ", ")[1]
		dateTimeString = strings.Split(dateTimeString, "Uhr")[0]
		return dateTimeString
	},
	titleSelector: ".views-field-title h2",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".views-field-title h2 a").Attr("href"); exists {
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

var heitereFahneCrawler = HTMLCrawler{
	venue:         Venue{ID: 2, Name: "Heitere Fahne", URL: "http://www.dieheiterefahne.ch/de/hauptnavigation/start/programm-31.html"},
	eventSelector: ".events .event",
	timeFormat:    "02.01.200615:04",
	getDateTimeString: func(eventSelection *goquery.Selection) string {
		rawDateTimeString := eventSelection.Find(".date + .time").Parent().Text()
		rawDateTimeString = StripSomeWhiteSpaces(rawDateTimeString)
		rawDateTimeString = strings.TrimSpace(rawDateTimeString)
		return rawDateTimeString[3:13] + rawDateTimeString[33:38]
	},
	titleSelector: ".alpha.omega.text .inner h2 a",
	linkBuilder: func(crawler *HTMLCrawler, eventSelection *goquery.Selection) string {
		if href, exists := eventSelection.Find(".alpha.omega.text .inner h2 a").Attr("href"); exists {
			// TODO FIXME the extracted href starts with /. Simply concatenating with the venue URl will not work.
			return fmt.Sprint(crawler.venue.URL, href)
		}
		return crawler.venue.URL // TODO set as default in Crawl if this function returns ""
	}}

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
	&heitereFahneCrawler}

// func main() {
// 	flag.Parse()
// 	events, err := kairoCrawler.Crawl()
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, ev := range events {
// 		fmt.Println(ev.Title, ev.DateTime, ev.URL)
// 	}
// }

// http://www.dieheiterefahne.ch/de/hauptnavigation/start/programm-31.html
// http://wartsaal-kaffee.ch/veranstaltungen/
// https://www.facebook.com/pg/CaffeBarSattler/events/t
// http://www.cafete.ch/
// http://www.cafemarta.ch/musik
// http://www.onobern.ch/programm-bersicht/
// http://www.schlachthaus.ch/spielplan/index.php
// http://dampfzentrale.ch/programm/
// http://www.bierhuebeli.ch/veranstaltungen/
// https://www.effinger.ch/events/
// https://www.facebook.com/pg/loescherbern/events/?ref=page_internal
// https://www.facebook.com/peterflamingobern/
// roessli, sous-le-pont,
