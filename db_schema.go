package wasgeit

import (
	"fmt"
	"io/ioutil"
)

func (store *Store) DropTables() error {
	if store.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}

	drop := readFile("sql/drop.sql")

	_, err := store.db.Exec(drop)

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) CreateTables() error {
	if store.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}

	schema := readFile("sql/create-schema.sql")
	venues := readFile("sql/insert-venues.sql")

	_, err := store.db.Exec(schema)

	if err != nil {
		return err
	}

	_, err = store.db.Exec(venues)

	if err != nil {
		return err
	}

	return nil
}

func readFile(filename string) string {
	schema, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(schema[:])
}
