package wasgeit

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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
	if st.db != nil {
		return st.db.Close()
	}

	return nil
}

func (st *Store) FindEvents(venue Venue) ([]Event, error) {
	rows, _ := st.db.Query("SELECT id, title, date, url FROM events where venue = ?", venue.ShortName)
	defer rows.Close()

	return mapRowsToEvents(rows)
}

func (st *Store) SaveEvent(ev Event) error {
	tx, err := st.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into events(title, date, url, venue) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(ev.Title, ev.DateTime, ev.URL, ev.Venue.ShortName)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to persists %v: %s", ev, err)
	}

	tx.Commit()
	return nil
}

func (st *Store) GetEvents() ([]Event, error) {
	rows, err := st.db.Query("SELECT id, title, date, url FROM events where date > DATE('now', '-1 day')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ct, _ := rows.ColumnTypes()
	for _, t := range ct {
		fmt.Printf("%v\n", t)
	}

	return mapRowsToEvents(rows)
}

func mapRowsToEvents(rows *sql.Rows) ([]Event, error) {
	var events []Event

	for rows.Next() {
		var ev Event
		// TODO map venues
		err := rows.Scan(&ev.ID, &ev.Title, &ev.DateTime, &ev.URL)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		events = append(events, ev)
	}

	return events, nil
}
