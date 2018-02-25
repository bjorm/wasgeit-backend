package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bjorm/wasgeit"
)

const tmpDataDir = "./tmp/"

func main() {
	crName := flag.String("name", "", "Name of crawler to run")
	flag.Parse()

	if *crName == "" {
		panic("Please specifiy a valid crawler")
	}

	st := wasgeit.Store{}
	st.Connect()

	wasgeit.RegisterAllHTMLCrawlers(&st)

	cr := wasgeit.GetCrawler(*crName).(*wasgeit.HTMLCrawler)

	var f *os.File
	filename := fmt.Sprintf("%s%s.html", tmpDataDir, cr.Name())

	if _, err := os.Stat(filename); err != nil {
		downloadSite(filename, cr)
	}

	f, err := os.Open(filename)
	panicOnError(err)
	defer f.Close()

	err = cr.LoadFrom(f)
	panicOnError(err)

	events, errors := cr.GetEvents()

	if len(events) == 0 {
		fmt.Println("No events returned.")
	}

	for _, ev := range events {
		fmt.Printf("title: %q\n", ev.Title)
		fmt.Printf("parsed time: %q\n", ev.DateTime)
		fmt.Printf("link: %q\n", ev.URL)
		fmt.Println()
	}

	for _, err := range errors {
		fmt.Println(err)
	}
}

func downloadSite(filename string, cr *wasgeit.HTMLCrawler) {
	resp, err := http.Get(cr.URL())
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
