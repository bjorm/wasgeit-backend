package wasgeit

import (
	"fmt"
)

// Venue describes a place where Events take place
type Venue struct {
	ID        int64
	ShortName string
	Name      string
	URL       string
	// TODO think about where this would be placed best
	IsSame    func(ev1, ev2 Event) bool
}

var Venues = []Venue{
	Venue{ID: 1, Name: "Cafe Kairo", ShortName: "kairo", URL: "http://www.cafe-kairo.ch/kultur", IsSame: hasSameUrl},
	Venue{ID: 2, Name: "Dachstock", ShortName: "dachstock", URL: "http://www.dachstock.ch", IsSame: hasSameUrl},
	Venue{ID: 3, Name: "Turnhalle", ShortName: "turnhalle", URL: "http://www.turnhalle.ch", IsSame: hasSameUrl},
	Venue{ID: 4, Name: "Brasserie Lorraine", ShortName: "brasserie-lorraine", URL: "http://brasserie-lorraine.ch/?post_type=tribe_events", IsSame: hasSameUrl},
	Venue{ID: 5, Name: "Kofmehl", ShortName: "kofmehl", URL: "http://www.kofmehl.net", IsSame: hasSameUrl},
	Venue{ID: 6, Name: "Kiff", ShortName: "kiff", URL: "http://www.kiff.ch", IsSame: hasSameUrl},
	Venue{ID: 7, Name: "Coq d'Or", ShortName: "coq-d-or", URL: "http://www.coq-d-or.ch/", IsSame: hasSameUrl},
	Venue{ID: 8, Name: "ISC", ShortName: "isc", URL: "http://www.isc-club.ch/", IsSame: hasSameUrl},
	Venue{ID: 9, Name: "Mahogany Hall", ShortName: "mahogany-hall", URL: "https://www.mahogany.ch/konzerte", IsSame: hasSameUrl},
	Venue{ID: 10, Name: "Heitere Fahne", ShortName: "heitere-fahne", URL: "http://www.dieheiterefahne.ch/de/hauptnavigation/start/programm-31.html", IsSame: hasSameUrl},
	Venue{ID: 11, Name: "ONO", ShortName: "ono", URL: "http://www.onobern.ch/programm-bersicht", IsSame: hasSameUrl},
	Venue{ID: 12, Name: "Cafe Marta", ShortName: "marta", URL: "http://www.cafemarta.ch/musik", IsSame: hasSameTitleAndDate},
	Venue{ID: 13, Name: "Bierhuebeli", ShortName: "bierhuebeli", URL: "http://www.bierhuebeli.ch/veranstaltungen/", IsSame: hasSameUrl},
	Venue{ID: 14, Name: "Dampfzentrale", ShortName: "dampfzentrale", URL: "http://dampfzentrale.ch/programm/", IsSame: hasSameUrl}}

func GetVenueOrPanic(shortName string) Venue {
	for _, venue := range Venues {
		if venue.ShortName == shortName {
			return venue
		}
	}
	panic(fmt.Sprintf("Failed to find venue with name %q", shortName))
}

func hasSameUrl(ev1, ev2 Event) bool {
	return ev1.URL == ev2.URL
}

func hasSameTitleAndDate(ev1, ev2 Event) bool {
	return ev1.Title == ev2.Title && ev1.DateTime.Equal(ev2.DateTime)
}