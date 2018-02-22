package wasgeit

import "fmt"

var schema = `
CREATE TABLE events (
	id INTEGER PRIMARY KEY, 
	title TEXT NOT NULL, 
	date DATETIME NOT NULL,
	url TEXT NOT NULL,
	venue TEXT NOT NULL,
	created DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX events_uq_title_date ON events(title, date);

CREATE TABLE venues (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	url	TEXT UNIQUE,
	name	TEXT UNIQUE,
	shortname	TEXT UNIQUE
);

CREATE TABLE logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	datetime  DATETIME,
	store_errors  TEXT,
	crawl_errors TEXT,
	updates	   	TEXT
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
