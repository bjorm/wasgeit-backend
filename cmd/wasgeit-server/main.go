package main

import (
	"flag"

	"github.com/golang/glog"
)

func main() {
	flag.Parse() // for glog

	dbErr := OpenDb()
	if dbErr != nil {
		panic(dbErr)
	}
	defer CloseDb()

	// dbErr = CreateTables()
	// if dbErr != nil {
	// 	panic(dbErr)
	// }

	for _, cr := range Crawlers {
		glog.V(1).Info(cr.Venue().Name, ":")
		
		events, err := cr.Crawl()

		if err != nil {
			glog.Infof("Error: %q", err)
		} else {
			for _, event := range events {
				storeErr := StoreEvent(event)
				if storeErr != nil {
					glog.Warningln(storeErr)
				}
			}
		}
	}

	glog.Flush()
}
