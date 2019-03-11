package wasgeit

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	st *Store
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

func (srv *Server) ServeAgenda(w http.ResponseWriter, r *http.Request) {
	events := srv.st.GetEventsYetToHappen()
	agenda := make(map[string][]interface{})

	for _, ev := range events {
		date := ev.DateTime.Format("2006-01-02")
		agenda[date] = append(agenda[date], from(ev))
	}
	b, err := json.Marshal(agenda)

	if err != nil {
		panic(err)
	}

	srv.addHeaders(w.Header())

	w.Write(b)
}

func (srv *Server) ServeNews(w http.ResponseWriter, r *http.Request) {
	events := srv.st.GetEventsAddedDuringLastWeek()
	news := make(map[string][]interface{})

	for _, ev := range events {
		date := ev.Created.Format("2006-01-02")
		news[date] = append(news[date], from(ev))
	}

	b, err := json.Marshal(news)

	if err != nil {
		panic(err)
	}

	srv.addHeaders(w.Header())

	w.Write(b)
}

func (srv *Server) addHeaders(h http.Header) {
	h.Add("ETag", srv.st.ReadValue(LastCrawlTimeKey))
	h.Add("Content-Type", "application/json;charset=utf-8")
}

func NewServer(st *Store) *Server {
	srv := Server{st: st}
	return &srv
}
