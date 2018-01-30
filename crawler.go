package wasgeit

import (
	"fmt"
	"time"
)

type RawEvent interface {
	Title() string
	DateTime() (time.Time, error)
	URL() string
}

type Crawler interface {
	URL() string
	Name() string
	Fetch() error
	GetEvents() ([]Event, []error)
	IsSame(ev1, ev2 Event) bool
}

func GetCrawler(name string) Crawler {
	if cr, exists := crawlers[name]; exists {
		return cr
	}
	return nil
}

func RegisterCrawler(name string, cr Crawler) {
	if _, exists := crawlers[name]; exists {
		panic(fmt.Sprintf("Crawler %q already registered.", name))
	} else {
		crawlers[name] = cr
	}
}

func GetCrawlers() []Crawler {
	var allcrawlers []Crawler
	for _, cr := range crawlers {
		allcrawlers = append(allcrawlers, cr)
	}
	return allcrawlers
}

var crawlers = make(map[string]Crawler)
