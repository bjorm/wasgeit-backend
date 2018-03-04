package wasgeit

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/goodsign/monday"
	"io"
)

type JsonCrawler struct {
	venue Venue
	body []byte
}

type GaskesselJson struct {
	Link string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
}

func (cr *JsonCrawler) URL() string {
	return cr.venue.URL
}

func (cr *JsonCrawler) Name() string {
	return cr.venue.ShortName
}

func (cr *JsonCrawler) Read(r io.ReadCloser) error {
	body, err := ioutil.ReadAll(r)
	defer r.Close()

	if err != nil {
		return err
	}
	cr.body = body

	return nil
}

func (cr *JsonCrawler) GetEvents() ([]Event, []error) {
	var jsonEvents []GaskesselJson
	var events []Event
	var errors []error
	err := json.Unmarshal(cr.body, &jsonEvents)

	if err != nil {
		return events, append(errors, err)
	}

	for _, jev := range jsonEvents {
		tokens := strings.Split(jev.Title.Rendered, " / ")
		if len(tokens) != 2 {
			errors = append(errors, fmt.Errorf("malformed title: %v", jev))
		}
		eventTime, err := monday.ParseInLocation("Mon 02.01.06", tokens[0], location, monday.LocaleDeDE)

		if err != nil {
			errors = append(errors, err)
		}

		title := strings.Replace(tokens[1], "<br>", "/ ", -1)
		events = append(events, Event{ Title: title, DateTime: eventTime, Venue: cr.venue, URL: jev.Link })
	}

	return events, errors
}

func (cr *JsonCrawler) IsSame(ev1, ev2 Event) bool {
	return ev1.URL == ev2.URL
}

func RegisterAllJsonCrawlers(st *Store) {
	RegisterCrawler("gaskessel", &JsonCrawler{ venue: st.GetVenue("gaskessel")})
}