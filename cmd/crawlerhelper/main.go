package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bjorm/wasgeit"
)

const tmpDataDir = "./tmp/"

func main() {
	crName := flag.String("name", "", "Name of crawler to run")
	config := wasgeit.GetConfiguration()
	wasgeit.ConfigureLogging(config.LogLevel)

	if *crName == "" {
		panic("Please specifiy a valid crawler")
	}

	st := wasgeit.Store{}
	st.Connect()

	browser := wasgeit.StartBrowser()
	defer browser.Close()

	wasgeit.RegisterAllHTMLCrawlers(&st)

	cr := wasgeit.GetCrawler(*crName)

	if cr == nil {
		panic(fmt.Sprintf("No crawler %q found", *crName))
	}

	var f *os.File
	filename := fmt.Sprintf("%s%s.%s", tmpDataDir, cr.Name(), inferExtension(cr))

	if _, err := os.Stat(filename); err != nil {
		downloadSite(filename, cr, browser)
	}

	bytes, err := ioutil.ReadFile(filename)
	panicOnError(err)
	defer f.Close()

	err = cr.Read(string(bytes))
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
func inferExtension(cr wasgeit.Crawler) string {
	switch cr.(type) {
	case *wasgeit.HTMLCrawler:
		return "html"
	default:
		return "txt"

	}
}

func downloadSite(filename string, cr wasgeit.Crawler, browser wasgeit.Browser) {
	body, err := browser.GetHtml(cr.URL())
	panicOnError(err)

	newLocalFile, err := os.Create(filename)
	panicOnError(err)
	defer newLocalFile.Close()

	_, err = newLocalFile.WriteString(body)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
