package crawler

import (
	"fmt"
	"os"
	"time"
	"github.com/huandu/facebook"
)

// Event describes an event taking place in a Venue
type Event struct {
	ID       int64
	Title    string
	DateTime time.Time
	URL      string
	Venue    *Venue
}

// Venue describes a place where Events take place
type Venue struct {
	ID   int64
	Name string
	URL  string
}

// Crawler describes  a crawler
type Crawler interface {
	Crawl() (events []*Event, err error)
}

type FbRoot struct {
	Events FbEvents
}

type FbEvents struct {
	Data []FbEvent
}

type FbEvent struct {
	ID string
	Name string
	StartTime string
	EndTime string
}

func (venue Venue) Crawl() (events []Event, err error) {
	apiKey, exists := os.LookupEnv("WASGEIT_FB_API_KEY")

	if !exists {
		return events, fmt.Errorf("no API key set")
	}

	res, _ := facebook.Get(fmt.Sprintf("/%s", venue.URL), facebook.Params{
		"fields":       "events",
		"time_filter": "upcoming",
		"access_token": apiKey,
	})

	var fbEvents FbRoot
	error := res.Decode(&fbEvents)

	if error != nil {
		fmt.Printf("Error: %s\n", error)
	} else {
		fmt.Printf("%v\n", fbEvents)
	}

	return events, nil
}

var Venues = []Venue{//Venue{Name: "Playground Lounge", URL: "playgroundlounge"},
	Venue{Name: "Brasserie Lorraine", URL: "106714019391041"},
	Venue{Name: "Dean Wake", URL: "113178252085227"},
	Venue{Name: "Kofmehl", URL: "tschieh"},
	Venue{Name: "Coq d'Or", URL: "Coq.d.Or"},
	Venue{Name: "KiFF", URL: "kiffaarau"},
	Venue{Name: "Fri-Son", URL: "frisonclub"},
	Venue{Name: "ISC", URL: "iscclub.bern"},
	Venue{Name: "Turnhalle", URL: "turnhalle"},
	Venue{Name: "Rössli", URL: "RossliBar"},
	Venue{Name: "Dachstock", URL: "dachstock"},
	Venue{Name: "Sous-le-Pont", URL: "slpreitschule"},
	Venue{Name: "Café Kairo", URL: "CafeKairo"}}
