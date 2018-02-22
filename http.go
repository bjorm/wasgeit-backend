package wasgeit

import (
	"encoding/json"
	"time"
	"net/http"
)

type Server struct {
	st *Store
}

type JsonEvent struct {
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Venue   Venue     `json:"venue"`
	Created time.Time `json:"created"`
}

func from(ev Event) JsonEvent {
	return JsonEvent{Title: ev.Title, URL: ev.URL, Venue: ev.Venue, Created: ev.Created}
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

	w.Write(b)
}

func NewServer(st *Store) *Server {
	srv := Server{st: st}
	return &srv
}
