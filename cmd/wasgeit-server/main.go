package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/bjorm/wasgeit"
)

func main() {
	flag.Parse() // for glog
	store := wasgeit.Store{}

	dbErr := store.Connect()
	if dbErr != nil {
		panic(dbErr)
	}
	defer store.Close()

	// dbErr = store.CreateTables()
	// if dbErr != nil {
	// panic(dbErr)
	// }

	for _, cr := range wasgeit.Crawlers {
		glog.V(1).Info(cr.Venue().Name, ":")
		
		events, err := cr.Crawl()

		if err != nil {
			glog.Infof("Error: %q", err)
		} else {
			for _, event := range events {
				storeErr := store.SaveEvent(event)
				if storeErr != nil {
					glog.Warningln(storeErr)
				}
			}
		}
	}

	glog.Flush()
}
