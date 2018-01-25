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
`

func (st *Store) CreateTables() error {
	if st.db == nil {
		return fmt.Errorf("Need to connect to DB first")
	}
	_, err := st.db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
