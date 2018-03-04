package main

import (
	"flag"
	"net/http"

	"github.com/bjorm/wasgeit"
	log "github.com/sirupsen/logrus"
)

func main() {
	resetDb := flag.Bool("setup-db", false, "Whether to create DB tables")
	flag.Parse()

	store := &wasgeit.Store{}
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

	wasgeit.RegisterAllHTMLCrawlers(store)
	wasgeit.RegisterAllJsonCrawlers(store)

	for _, cr := range wasgeit.GetCrawlers() {
		log.Info(cr.Name())

		resp, err := http.Get(cr.URL())
		if err != nil {
			log.Errorf("Fetching failed: %s", err)
			continue
		}

		err = cr.Read(resp.Body)

		if err != nil {
			log.Errorf("Reading failed: %s", err)
			continue
		}

		newEvents, crawlErrors := cr.GetEvents()

		if len(newEvents) == 0 {
			log.Errorf("Crawler %q returned no events", cr.Name())
			continue
		}

		// TODO use channel and goroutines for this

		existingEvents := store.FindEvents(cr.Name())

		if len(existingEvents) == 0 {
			log.Warnf("No existing events found", cr.Name())
		}

		cs := wasgeit.DedupeAndTrackChanges(existingEvents, newEvents, cr)
		var storeErrors []error

		for _, event := range cs.New {
			storeErr := store.SaveEvent(event)

			if storeErr != nil {
				storeErrors = append(storeErrors, storeErr)
			}
		}

		// TODO do updates

		log.Infof("Crawl errors: %s", crawlErrors)
		log.Infof("Store errors: %s", storeErrors)
		log.Infof("Updates: %+v", cs.Updates)
		log.Infof("New events stored: %d", len(cs.New)-len(storeErrors))
	}

	server := wasgeit.NewServer(store)
	http.HandleFunc("/agenda", server.ServeAgenda)
	http.HandleFunc("/news", server.ServeNews)

	log.Info("Serving..")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
