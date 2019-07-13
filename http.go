package wasgeit

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	store *Store
}

type JsonEvent struct {
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	DateTime time.Time `json:"datetime"`
	Venue    Venue     `json:"venue"`
	Created  time.Time `json:"created"`
}

func from(ev Event) JsonEvent {
	return JsonEvent{Title: ev.Title, URL: ev.URL, DateTime: ev.DateTime, Venue: ev.Venue, Created: ev.Created}
}

func (server *Server) ServeAgenda(w http.ResponseWriter, r *http.Request) {
	events := server.store.GetEventsYetToHappen()
	agenda := make(map[string][]interface{})

	for _, ev := range events {
		date := ev.DateTime.Format("2006-01-02")
		agenda[date] = append(agenda[date], from(ev))
	}
	b, err := json.Marshal(agenda)

	if err != nil {
		panic(err)
	}

	server.setContentType(w.Header())
	server.setEtag(w.Header())

	w.Write(b)
}

func (server *Server) ServeNews(w http.ResponseWriter, r *http.Request) {
	events := server.store.GetEventsAddedDuringLastWeek()
	news := make(map[string][]interface{})

	for _, ev := range events {
		date := ev.Created.Format("2006-01-02")
		news[date] = append(news[date], from(ev))
	}

	b, err := json.Marshal(news)

	if err != nil {
		panic(err)
	}

	server.setContentType(w.Header())
	server.setEtag(w.Header())

	w.Write(b)
}

func (server *Server) ServeFestivals(w http.ResponseWriter, r *http.Request) {
	festivals, err := server.store.GetCurrentFestivals()

	if err != nil {
		log.Error(err)
	}

	b, err := json.Marshal(festivals)

	if err != nil {
		panic(err)
	}

	server.setContentType(w.Header())

	w.Write(b)
}

func (server *Server) setContentType(h http.Header) {
	h.Add("Content-Type", "application/json;charset=utf-8")
}

func (server *Server) setEtag(h http.Header) {
	h.Add("ETag", server.store.ReadValue(LastCrawlTimeKey))
}

func NewServer(st *Store) *Server {
	srv := Server{store: st}
	return &srv
}
