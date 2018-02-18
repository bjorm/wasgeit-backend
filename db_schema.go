package wasgeit

import "fmt"

var schema = `
CREATE TABLE events (
	id INTEGER PRIMARY KEY, 
	title TEXT, 
	date DATETIME, 
	url TEXT,
	venue TEXT
);
CREATE UNIQUE INDEX events_uq_title_date ON events(title, date);
CREATE TRIGGER timestamp_created AFTER INSERT ON events BEGIN UPDATE events SET created = DATETIME('now') WHERE id = NEW.id; END;

CREATE TABLE venues (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	url	TEXT UNIQUE,
	name	TEXT UNIQUE,
	shortname	TEXT UNIQUE
);

CREATE TABLE logs (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	datetime	TEXT,
	store_errors	TEXT,
	crawl_errors	TEXT,
	updates	TEXT
);
`

func (st *Store) CreateTables() error {
	if st.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}
	_, err := st.db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
