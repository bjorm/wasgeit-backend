package main

import (
	"time"

	"github.com/bjorm/wasgeit"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := wasgeit.GetConfiguration()

	wasgeit.ConfigureLogging(config.LogLevel)

	store := &wasgeit.Store{}
	dbErr := store.Connect()

	if dbErr != nil {
		panic(dbErr)
	}
	defer store.Close()

	if config.DropDb {
		log.Info("Dropping DB..")
		dbErr = store.DropTables()
		if dbErr != nil {
			panic(dbErr)
		}
	}

	if config.SetupDb {
		log.Info("Setting up DB..")
		dbErr = store.CreateTables()
		if dbErr != nil {
			panic(dbErr)
		}
	}

	wasgeit.RegisterAllHTMLCrawlers(store)

	browser, err := wasgeit.StartBrowser(config.ChromiumUrl)

	if err != nil {
		panic(err)
	}

	defer browser.Close()

	for _, cr := range wasgeit.GetCrawlers() {
		log.Info(cr.Name())

		body, err := browser.GetHtml(cr.URL())

		log.Debug("Got site body from browser")

		if err != nil {
			log.Errorf("Fetching failed: %s", err)
			continue
		}

		err = cr.Read(body)

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
