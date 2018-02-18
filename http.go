package wasgeit

import (
	"encoding/json"
	"time"
	"net/http"
)

type Server struct {
	st *Store
}

func (srv *Server) ServeAgenda(w http.ResponseWriter, r *http.Request) {
	events := srv.st.GetEventsYetToHappen()
	agenda := make(map[string][]interface{})

	now := time.Now()
	todayMorning := time.Date(now.Year(), now.Month(), now.Day(), 0, 0 ,0, 0, location)

	for _, ev := range events {
		if  ev.DateTime.After(todayMorning) {
			date := ev.DateTime.Format("2006-01-02")
			// oh yeah
			ev2 := struct {
				Title string `json:"title"`
				URL string `json:"url"`
				Venue Venue `json:"venue"`
			}{ Title: ev.Title, URL: ev.URL, Venue: ev.Venue }
			agenda[date] = append(agenda[date], ev2)
		}
	}
	b, err := json.Marshal(agenda)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func NewServer(st *Store) *Server {
	srv := Server{st: st}
	return &srv
}
