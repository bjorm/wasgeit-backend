package wasgeit

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("wasgeit-server")

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
		return fmt.Errorf("Failed to persists %v: %s", ev, err)
	}

	tx.Commit()
	return nil
}

func (st *Store) GetEvents() ([]Event, error) {
	var events []Event
	rows, err := st.db.Query("SELECT id, title, date, url FROM events where date > DATE('now', '-1 day')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ct, _ := rows.ColumnTypes()
	for _, t := range ct {
		fmt.Printf("%v\n", t)
	}

	for rows.Next() {
		var ev Event
		err = rows.Scan(&ev.ID, &ev.Title, &ev.DateTime, &ev.URL)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		events = append(events, ev)
	}
	return events, nil
}
