package wasgeit

import (
	"fmt"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

func (st *Store) CreateTables() error {
	if st.db == nil {
		return fmt.Errorf("need to connect to DB first")
	}

	drop := readFile("sql/drop.sql")
	schema := readFile("sql/create-schema.sql")
	venues := readFile("sql/insert-venues.sql")

	logrus.Info("Dropping tables, etc.. ")
	_, err := st.db.Exec(drop)

	if err != nil {
		return err
	}

	logrus.Info("Creating schema..")
	_, err = st.db.Exec(schema)

	if err != nil {
		return err
	}

	logrus.Info("Inserting venues..")
	_, err = st.db.Exec(venues)

	if err != nil {
		return err
	}

	logrus.Infoln("Done")

	return nil
}

func readFile(filename string) string {
	schema, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(schema[:])
}
