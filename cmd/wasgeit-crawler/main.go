package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/bjorm/wasgeit"
	log "github.com/sirupsen/logrus"
)

func main() {
	dropDb := flag.Bool("drop-db", false, "Whether to drop DB")
	setupDb := flag.Bool("setup-db", false, "Whether to create DB tables")
	flag.Parse()

	store := &wasgeit.Store{}
	dbErr := store.Connect()

	if dbErr != nil {
		panic(dbErr)
	}
	defer store.Close()

	if *dropDb {
		log.Info("Dropping DB..")
		dbErr = store.DropTables()
		if dbErr != nil {
			panic(dbErr)
		}
	}

	if *setupDb {
		log.Info("Setting up DB..")
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
			log.Warnf("No existing events found")
		}

		cs := wasgeit.DedupeAndTrackChanges(existingEvents, newEvents, cr)
		var storeErrors []error

		for _, update := range cs.Updates {
			for _, field := range update.ChangedFields {
				var newValue, oldValue interface{}
				switch field {
				case "title":
					newValue = update.UpdatedEv.Title
					oldValue = update.ExistingEv.Title
					break
				case "date":
					newValue = update.UpdatedEv.DateTime
					oldValue = update.ExistingEv.DateTime
				default:
					panic("Update not implemented.")
				}
				store.UpdateEvent(update.ExistingEv.ID, field, newValue)
				store.LogUpdate(update.ExistingEv.ID, field, oldValue, newValue)
			}
		}

		for _, event := range cs.New {
			storeErr := store.SaveEvent(event)

			if storeErr != nil {
				storeErrors = append(storeErrors, storeErr)
			}
		}

		for _, err := range storeErrors {
			store.LogError(cr, err)
		}

		log.Infof("Crawl errors: %d", len(crawlErrors))
		log.Infof("Store errors: %d", len(storeErrors))
		log.Infof("Updates: %d", len(cs.Updates))
		log.Infof("New events stored: %d", len(cs.New)-len(storeErrors))
	}

	store.UpdateValue(wasgeit.LastCrawlTimeKey, time.Now().Format(time.RFC3339))
}
