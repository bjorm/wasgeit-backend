package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
)

func main() {
	flag.Parse() // demanded by glog

	for _, cr := range Crawlers {
		fmt.Printf("%s:\n", cr.Venue().Name)

		events, err := cr.Crawl()

		if err != nil {
			glog.Infof("Error: %q", err)
		} else {
			for _, event := range events {
				fmt.Printf("%s: %s\n%s\n", event.DateTime.Format("02 Jan"), event.Title, event.URL)
			}
			fmt.Println("")
		}
	}

	glog.Flush()
}
