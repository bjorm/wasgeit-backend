package wasgeit

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

const schemaVersion = 1

type Store struct {
	db *sql.DB
}

func (store *Store) Connect() error {
	var err error

	store.db, err = sql.Open("sqlite3", "db/wasgeit.db")

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) Close() error {
	if store.db != nil {
		return store.db.Close()
	}

	return nil
}

func (store *Store) FindVenue(shortName string) (Venue, error) {
	row := store.db.QueryRow("SELECT id, name, shortname, url FROM venues WHERE shortname = ?", shortName)
	var v Venue
	err := row.Scan(&v.ID, &v.Name, &v.ShortName, &v.URL)

	if err == sql.ErrNoRows {
		return Venue{}, fmt.Errorf("could not find venue %q", shortName)
	} else if err != nil {
		return Venue{}, fmt.Errorf("querying venues failed: %q", err)
	}
	return v, nil
}

func (store *Store) GetVenue(shortName string) Venue {
	venue, err := store.FindVenue(shortName)

	if err != nil {
		panic(err)
	}
	return venue
}

func (store *Store) FindEvents(crawlerName string) []Event {
	rows, err := store.db.Query(`SELECT 
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

func (store *Store) SaveEvent(ev Event) error {
	return store.inTransaction(`insert into events(title, date, url, venue) values(?, ?, ?, ?)`, func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(ev.Title, ev.DateTime, ev.URL, ev.Venue.ShortName)
	}, func(err error) error {
		return fmt.Errorf("failed to persists event %v: %s", ev, err)
	})
}

func (store *Store) GetEventsYetToHappen() []Event {
	rows, err := store.db.Query(`SELECT 
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

func (store *Store) GetEventsAddedDuringLastWeek() []Event {
	rows, err := store.db.Query(`SELECT 
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

func (store *Store) UpdateEvent(id int64, fieldName string, value interface{}) {
	if fieldName != "title" && fieldName != "date" {
		panic(fmt.Sprintf("Unknown column provided for update: %q", fieldName))
	}

	updateQuery := fmt.Sprintf("UPDATE events SET %s = ? WHERE id = ?", fieldName)

	err := store.inTransaction(updateQuery, func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(value, id)
	}, func(err error) error {
		return fmt.Errorf("failed to update %q in event %d to %q because of: %s", fieldName, id, value, err)
	})

	if err != nil {
		panic(err)
	}
}

func (store *Store) LogUpdate(eventId int64, fieldName string, oldValue interface{}, newValue interface{}) {
	err := store.inTransaction(`INSERT INTO updates (event_id, field, old, new) VALUES (?, ?, ?, ?)`, func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(eventId, fieldName, oldValue, newValue)
	}, func(err error) error {
		return fmt.Errorf("failed to update %q in event %d, oldValue=%s, newValue=%s", fieldName, eventId, oldValue, newValue)
	})

	if err != nil {
		panic(err)
	}
}

func (store *Store) LogError(cr Crawler, errToLog error) {
	err := store.inTransaction(`INSERT INTO errors (crawler, msg) VALUES (?, ?)`, func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(cr.Name(), errToLog.Error())
	}, func(err error) error {
		return fmt.Errorf("failed to store error %q for %q", err, cr.Name())
	})

	if err != nil {
		panic(err)
	}
}

func (store *Store) UpdateValue(key string, newValue string) {
	err := store.inTransaction("INSERT OR REPLACE INTO keyvalue (key, value) VALUES (?, ?)", func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(key, newValue)
	}, func(err error) error {
		return fmt.Errorf("failed to set value of %q to %q", key, newValue)
	})

	if err != nil {
		panic(err)
	}
}

func (store *Store) ReadValue(key string) string {
	rows, err := store.db.Query("SELECT value FROM keyvalue WHERE key=(?)", key)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var value string

	if rows.Next() {
		err = rows.Scan(&value)
		if err != nil {
			panic(err)
		}
	}

	if rows.Next() {
		panic(fmt.Sprintf("Lookup of key %q in key value table returned more than one row", key))
	}

	return value
}

func (store *Store) inTransaction(query string, exec func(stmt *sql.Stmt) (sql.Result, error), createError func(err error) error) error {
	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = exec(stmt)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return createError(fmt.Errorf("Both rollack (%v) and statement (%v) failed ", rollbackErr, err))
		}
		return createError(err)
	}

	err = tx.Commit()

	if err != nil {
		return createError(err)
	}

	return nil
}

func (store *Store) GetCurrentFestivals() ([]Festival, error) {
	festivals := make([]Festival, 0)

	rows, err := store.db.Query("SELECT id, url, location, name, date_start, date_end FROM venues WHERE placement = 'what-else'")

	if err != nil {
		return festivals, fmt.Errorf("error when getting festivals: %v", err)
	}

	for rows.Next() {
		var festival Festival
		err := rows.Scan(&festival.Id, &festival.Url, &festival.Location, &festival.Title, &festival.DateStart, &festival.DateEnd)

		if err != nil {
			return festivals, fmt.Errorf("error when getting festivals: %v", err)
		}

		openingTimes, err := store.getOpeningTimes(festival.Id)

		if err != nil {
			return festivals, fmt.Errorf("could not get opening times for festival %q: %v", festival.Title, err)
		}

		festival.OpeningTimes = openingTimes

		festivals = append(festivals, festival)
	}

	return festivals, nil
}

func (store *Store) getOpeningTimes(venueId int64) ([]OpeningTime, error) {
	rows, err := store.db.Query("SELECT days, time_start, time_end FROM opening_times WHERE venue_id = ?", venueId)

	var openingTimes []OpeningTime

	log.Tracef("Fetching opening times for festival with id=%d", venueId)

	if err != nil {
		return openingTimes, fmt.Errorf("error when getting opening times: %v", err)
	}

	for rows.Next() {
		var openingTime OpeningTime
		err := rows.Scan(&openingTime.Days, &openingTime.Start, &openingTime.End)

		if err != nil {
			return openingTimes, fmt.Errorf("error when getting opening times: %v", err)
		}

		openingTimes = append(openingTimes, openingTime)
	}

	log.Tracef("Found %d opening times for festival with id=%d", len(openingTimes), venueId)

	return openingTimes, nil
}
