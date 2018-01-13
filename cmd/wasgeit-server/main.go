package main

import (
	"flag"
	"net/http"

	"github.com/bjorm/wasgeit"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("wasgeit")

func main() {
	resetDb := flag.Bool("setup-db", false, "Whether to create DB tables")
	flag.Parse()

	store := wasgeit.Store{}

	dbErr := store.Connect()
	if dbErr != nil {
		panic(dbErr)
	}
	defer store.Close()

	if *resetDb {
		log.Info("Setting up DB tables..")
		dbErr = store.CreateTables()
		if dbErr != nil {
			panic(dbErr)
		}
	}

	for _, cr := range wasgeit.Crawlers {
		log.Info(cr.Venue().Name)

		doc, err := cr.Get()
		events, crawlErrors := cr.Crawl(doc)

		if err != nil {
			log.Infof("Getting document for %q failed: %s", cr.Venue().Name, err)
			break
		}

		var storeErrors []error

		for _, event := range events {
			storeErr := store.SaveEvent(event)

			if storeErr != nil {
				storeErrors = append(storeErrors, storeErr)
			}
		}

		log.Infof("Crawl errors: %s", crawlErrors)
		log.Infof("Store errors: %s", storeErrors)
		log.Infof("Crawled and stored %d events successfully", len(events)-len(storeErrors))
	}

	http.Handle("/events", wasgeit.NewServer(&store))
	log.Info("Serving..")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
