package wasgeit

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Browser struct {
	ctxt   context.Context
	cancel context.CancelFunc
}

type WebsocketTargetJson struct {
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
	Url                  string `json:"url"`
}

func StartBrowser(chromiumHost string) Browser {
	hostUrl, err := url.Parse(chromiumHost)
	if err != nil {
		panic(err)
	}

	fullUrl, err := hostUrl.Parse("/json")

	if err != nil {
		panic(err)
	}

	wsUrl := getWsUrl(fullUrl)

	log.Debug("Got websocket URL ", wsUrl, " from ", fullUrl)
	ctxt, cancel := chromedp.NewRemoteAllocator(context.Background(), wsUrl)
	log.Debug("Connected to remote chromium")

	return Browser{ctxt: ctxt, cancel: cancel}
}

func getWsUrl(chromiumRemoteUrl *url.URL) string {
	resp, err := http.Get(chromiumRemoteUrl.String())

	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	targets := make([]WebsocketTargetJson, 0)
	err = json.Unmarshal(bytes, &targets)

	if err != nil {
		panic(err)
	}

	for _, target := range targets {
		if target.Url == "about:blank" { // we assume this is the default page
			return target.WebSocketDebuggerUrl
		}
	}

	panic("No default page found")
}

func (b *Browser) GetHtml(url string) (string, error) {
	log.Debug("Opening new tab for ", url)

	ctxt, cancel := chromedp.NewContext(b.ctxt) // create new tab
	defer cancel()                              // close tab

	log.Trace("Created")

	log.Trace("Running tasks..")

	var body string
	if err := chromedp.Run(ctxt,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		return body, err
	}

	return body, nil
}

func (b *Browser) Close() {
	log.Debug("Closing chromium")
	b.cancel()
}
