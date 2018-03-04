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
		return Venue{}, fmt.Errorf("could not find venue %q", shortName)
	} else if err != nil {
		return Venue{}, fmt.Errorf("querying venues failed: %q", err)
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

func (st *Store) FindEvents(crawlerName string) []Event {
	rows, err := st.db.Query(`SELECT 
		events.id, 
		events.title, 
		events.date, 
		events.url, 
		events.created,
		venues.id,
		venues.name,
		venues.shortname,
		venues.url
		FROM events 
		JOIN venues ON venues.shortname = events.venue
		WHERE venue = ?`,
		crawlerName)
	if err != nil {
		panic(err)
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
		return fmt.Errorf("failed to persists %v: %s", ev, err)
	}

	tx.Commit()
	return nil
}

func (st *Store) GetEventsYetToHappen() []Event {
	rows, err := st.db.Query(`SELECT 
								events.id, 
								events.title, 
								events.date, 
								events.url, 
								events.created,
								venues.id,
								venues.name,
								venues.shortname,
								venues.url
								FROM events 
								JOIN venues ON venues.shortname = events.venue 
								WHERE date(date) > DATE('now', '-1 day')`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return mapRowsToEvents(rows)
}

func (st *Store) GetEventsAddedDuringLastWeek() []Event {
	rows, err := st.db.Query(`SELECT 
								events.id, 
								events.title, 
								events.date, 
								events.url,
								events.created, 
								venues.id,
								venues.name,
								venues.shortname,
								venues.url
								FROM events 
								JOIN venues ON venues.shortname = events.venue
								WHERE date(created) > DATE('now', '-7 day') ORDER BY created DESC`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return mapRowsToEvents(rows)
}

func mapRowsToEvents(rows *sql.Rows) []Event {
	var events []Event

	for rows.Next() {
		var ev Event
		err := rows.Scan(&ev.ID, &ev.Title, &ev.DateTime, &ev.URL, &ev.Created, &ev.Venue.ID, &ev.Venue.Name, &ev.Venue.ShortName, &ev.Venue.URL)

		if err != nil {
			panic(err)
		}
		events = append(events, ev)
	}

	return events
}

func (st *Store) UpdateEvent(id int64, columnName string, value interface{}) {
	if columnName != "title" && columnName != "date" {
		panic(fmt.Sprintf("Unknown column provided for update: %q", columnName))
	}

	updateQuery := fmt.Sprintf("UPDATE events SET %s = ? WHERE id = ?", columnName)
	_, err := st.db.Exec(updateQuery, columnName, value, id)

	if err != nil {
		panic(err)
	}
}

func (st *Store) LogUpdate(eventId int64, fieldName string, oldValue interface{}, newValue interface{}) {
	_, err := st.db.Exec(
		`INSERT INTO updates (event_id, field, old, new) VALUES (?, ?, ?, ?)`,
		eventId, fieldName, oldValue, newValue)

	if err != nil {
		panic(err)
	}
}

func (st *Store) LogError(cr Crawler, errToLog error) {
	_, err := st.db.Exec(`INSERT INTO errors (crawler, msg) VALUES (?, ?)`, cr.Name(), errToLog.Error())

	if err != nil {
		panic(err)
	}
}