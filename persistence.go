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

func (st *Store) FindVenue(shortName string) (Venue, error) {
	row := st.db.QueryRow("SELECT id, name, shortname, url FROM venues WHERE shortname = ?", shortName)
	var v Venue
	err := row.Scan(&v.ID, &v.Name, &v.ShortName, &v.URL)

	if err == sql.ErrNoRows {
		return Venue{}, fmt.Errorf("Could not find venue %q", shortName)
	} else if err != nil {
		return Venue{}, fmt.Errorf("Error when querying venues: %q", err)
	}
	return v, nil
}

func (st *Store) GetVenue(shortName string) Venue {
	venue, err := st.FindVenue(shortName)

	if err != nil {
		panic(err)
	}
	return venue
}

func (st *Store) FindEvents(crawlerName string) ([]Event, error) {
	rows, err := st.db.Query("SELECT id, title, date, url FROM events WHERE venue = ?", crawlerName)
	if err != nil {
		return []Event{}, err
	}
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

	return mapRowsToEvents(rows)
}

func mapRowsToEvents(rows *sql.Rows) ([]Event, error) {
	var events []Event

	for rows.Next() {
		var ev Event
		// TODO map venues
		err := rows.Scan(&ev.ID, &ev.Title, &ev.DateTime, &ev.URL)

		if err != nil {
			return nil, err
		}
		events = append(events, ev)
	}

	return events, nil
}
