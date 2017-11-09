package main

import (
	"flag"
	"github.com/bjorm/wasgeit/crawler"
	"github.com/golang/glog"
)

func main() {
	flag.Parse() // demanded by glog

	for _, cr := range crawler.Crawlers {
		events, err := cr.Crawl()
		if err != nil {
			glog.Infof("Error: %q", err)
		} else {
			glog.Infof("%#v", events)
		}
	}

	glog.Flush()
}

