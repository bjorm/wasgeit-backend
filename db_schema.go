package wasgeit

import (
	"fmt"
	"io/ioutil"
)

func (st *Store) DropTables() error {
	if st.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}

	drop := readFile("sql/drop.sql")

	_, err := st.db.Exec(drop)

	if err != nil {
		return err
	}

	return nil
}

func (st *Store) CreateTables() error {
	if st.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}

	schema := readFile("sql/create-schema.sql")
	venues := readFile("sql/insert-venues.sql")

	_, err := st.db.Exec(schema)

	if err != nil {
		return err
	}

	_, err = st.db.Exec(venues)

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
