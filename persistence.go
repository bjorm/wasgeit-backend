package wasgeit

import (
	"database/sql"
	"github.com/op/go-logging"
	_ "github.com/mattn/go-sqlite3"
)

var log = logging.MustGetLogger("persistence")

const schemaVersion = 1

type Store struct {
	db *sql.DB
}

func (st *Store) Connect() error {
	var err error

	st.db, err = sql.Open("sqlite3", "db/wasgeit.db")

	if err != nil {
		return err
	}

	return nil
}

func (st *Store) Close() error {
	log.Info("Closing db")

	if st.db != nil {
		return st.db.Close()
	}

	return nil
}

func (st *Store) SaveEvent(ev Event) error {
	tx, err := st.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into events(title, date, url) values(?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(ev.Title, ev.DateTime, ev.URL)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (st *Store) GetEvents() []Event {
	st.db.Query("SELECT * FROM events where date > DATE('now', '-1 day')")
	return nil
}