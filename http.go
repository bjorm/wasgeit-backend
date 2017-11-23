package wasgeit

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type Server struct {
	st *Store
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, _ := srv.st.GetEvents()
	b, err := json.Marshal(events)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(b)
}

func NewServer(st *Store) *Server {
	srv := Server{st: st}
	return &srv
}