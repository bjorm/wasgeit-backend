package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/bjorm/wasgeit"
)

const tmpDataDir = "./tmp/"

func main() {
	crName := flag.String("name", "", "Name of crawler to run")
	flag.Parse()

	if *crName == "" {
		panic("Please specifiy a valid crawler")
	}

	if cr, exists := wasgeit.HTMLCrawlers[*crName]; exists {
		var f *os.File
		filename := fmt.Sprintf("%s%s.html", tmpDataDir, cr.Venue().ShortName)

		if _, err := os.Stat(filename); err != nil {
			downloadSite(filename, &cr)
		}

		f, err := os.Open(filename)
		panicOnError(err)
		defer f.Close()

		doc, err := goquery.NewDocumentFromReader(f)
		panicOnError(err)

		doc.Find(cr.EventSelector).Each(func(_ int, firstEv *goquery.Selection) {
			dt := cr.GetDateTimeString(firstEv)
			pdt, _ := cr.GetEventTime(firstEv)
			fmt.Printf("time string: %q\n", dt)
			fmt.Printf("parsed time: %q\n", pdt)
			fmt.Printf("link: %q\n", cr.LinkBuilder(&cr, firstEv))
			fmt.Printf("title: %q\n", firstEv.Find(cr.TitleSelector).Text())
			fmt.Println()
		})
	} else {
		panic("Crawler not found")
	}
}

func downloadSite(filename string, cr *wasgeit.HTMLCrawler) {
	resp, err := http.Get(cr.Venue().URL)
	panicOnError(err)

	body, err := ioutil.ReadAll(resp.Body)
	panicOnError(err)

	newLocalFile, err := os.Create(filename)
	panicOnError(err)
	defer newLocalFile.Close()

	_, err = newLocalFile.Write(body)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
