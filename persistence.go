package main

import (
	"database/sql"
	"fmt"

	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDb() error {
	var err error

	glog.Infof("Opening db")
	db, err = sql.Open("sqlite3", "db/wasgeit.db")

	if err != nil {
		return err
	}
	return nil
}

func CloseDb() error {
	glog.Infof("Closing db")
	if db != nil {
		return db.Close()
	}
	return nil
}

func StoreEvent(ev Event) error {
	tx, err := db.Begin()
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

var schema = `
CREATE TABLE venues (id INTEGER PRIMARY KEY, name TEXT UNIQUE, url TEXT);
CREATE TABLE events (id INTEGER PRIMARY KEY, title TEXT, date TEXT, url TEXT UNIQUE);
`

func CreateTables() error {
	if db == nil {
		return fmt.Errorf("Need to connect to DB first")
	}
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
