package wasgeit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
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

func StartBrowser(chromiumHost string) (Browser, error) {
	hostUrl, err := url.Parse(chromiumHost)
	if err != nil {
		panic(err)
	}

	fullUrl, err := hostUrl.Parse("/json")

	if err != nil {
		panic(err)
	}

	wsUrl, err := getWsUrl(fullUrl)

	if err != nil {
		return Browser{}, err
	}

	log.Debugf("Got websocket URL %s from %q", wsUrl, fullUrl)
	ctxt, cancel := chromedp.NewRemoteAllocator(context.Background(), wsUrl)
	log.Debugf("Connected to remote chromium at %s", wsUrl)

	return Browser{ctxt: ctxt, cancel: cancel}, nil
}

// replaceResolvedHostnameIfNeeded tries to resolve the hostname in the given URL in order to work around mitigations
// introduced because of https://bugs.chromium.org/p/chromium/issues/detail?id=813540
//
// Chromium will return the error "Host header is specified and is not an IP address or localhost" if the
// DevTools are connected to via a hostname other than localhost.
func replaceResolvedHostnameIfNeeded(url *url.URL) error {
	if url.Hostname() == "localhost" {
		return nil
	}

	addresses, err := net.LookupHost(url.Hostname())

	if err != nil {
		return err
	}

	if len(addresses) != 1 {
		log.Warnf("Warning: Host lookup of %q returned %v, picking first address\n", url.Hostname(), addresses)
	}

	port := url.Port()
	url.Host = fmt.Sprintf("%s:%s", addresses[0], port)

	log.Infof("Replaced hostname with %q\n", addresses[0])

	return nil
}

func getWsUrl(chromiumRemoteUrl *url.URL) (string, error) {
	err := replaceResolvedHostnameIfNeeded(chromiumRemoteUrl)

	if err != nil {
		return "", err
	}

	resp, err := http.Get(chromiumRemoteUrl.String())

	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	err = resp.Body.Close()

	if err != nil {
		return "", err
	}

	targets := make([]WebsocketTargetJson, 0)
	err = json.Unmarshal(bytes, &targets)

	if err != nil {
		return "", fmt.Errorf("could not parse %q as JSON, error was: %q", string(bytes), err)
	}

	for _, target := range targets {
		if target.Url == "about:blank" { // we assume this is the default page
			return target.WebSocketDebuggerUrl, nil
		}
	}

	return "", fmt.Errorf("no default page found on %q", chromiumRemoteUrl.String())
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
