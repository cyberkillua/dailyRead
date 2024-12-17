package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Atom struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title string `xml:"title"`
	Link  struct {
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Description string `xml:"summary"`
	PublishedAt string `xml:"published"`
}

type RSS struct {
	XMLName xml.Name       `xml:"rss"`
	Version string         `xml:"version,attr"`
	Channel GenericChannel `xml:"channel,omitempty"`
}
type GenericChannel struct {
	Title       string    `xml:"title,omitempty"`
	Link        string    `xml:"link,omitempty"`
	Description string    `xml:"description,omitempty"`
	Language    string    `xml:"language,omitempty"`
	Items       []RSSItem `xml:"item,omitempty"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
}

func convertAtomToRSSItems(entries []AtomEntry) []RSSItem {
	rssItems := make([]RSSItem, len(entries))
	for i, entry := range entries {
		rssItems[i] = RSSItem{
			Title:       entry.Title,
			Link:        entry.Link.Href,
			Description: entry.Description,
			PubDate:     entry.PublishedAt, // Atom dates are ISO 8601; RSS uses RFC 1123
		}
	}
	return rssItems
}

func urlToRSS(url string) (RSS, error) {
	// Create a custom transport to handle redirects more explicitly
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	httpClient := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		// Explicitly handle redirects to get more information
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("stopped after 5 redirects")
			}

			if len(via) > 0 {
				log.Printf("Redirect from %s to %s", via[len(via)-1].URL, req.URL)
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RSS{}, fmt.Errorf("failed to create request: %w", err)
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Feedfetcher-Google; (+http://www.google.com/feedfetcher.html)",
	}

	req.Header.Set("User-Agent", userAgents[0])
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := httpClient.Do(req)
	if err != nil {
		return RSS{}, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	// Extensive logging for debugging
	log.Printf("Request URL: %s", url)
	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Content-Type: %s", resp.Header.Get("Content-Type"))

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSS{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug: log first 500 characters
	// log.Printf("Response Body (first 500 chars): %s", string(data[:min(len(data), 500)]))

	processedData := preprocessXML(data)

	var root struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(processedData, &root); err != nil {
		return RSS{}, fmt.Errorf("failed to parse XML: %w", err)
	}

	switch strings.ToLower(root.XMLName.Local) {
	case "rss":
		var rssFeed RSS
		if err := xml.Unmarshal(processedData, &rssFeed); err != nil {
			return RSS{}, fmt.Errorf("failed to parse RSS feed: %w", err)
		}
		return rssFeed, nil
	case "feed":
		var atomFeed Atom
		if err := xml.Unmarshal(processedData, &atomFeed); err != nil {
			return RSS{}, fmt.Errorf("failed to parse Atom feed: %w", err)
		}
		rssFeed := RSS{
			Channel: GenericChannel{
				Title: atomFeed.Title,
				Items: convertAtomToRSSItems(atomFeed.Entries),
			},
		}
		return rssFeed, nil
	default:
		return RSS{}, fmt.Errorf("unknown feed format: %s", root.XMLName.Local)
	}
}

func preprocessXML(data []byte) []byte {

	data = bytes.TrimSpace(data)
	data = regexp.MustCompile(`<!--.*?-->`).ReplaceAll(data, []byte{})

	// Replace problematic character entities
	replacements := []struct{ old, new []byte }{
		{[]byte("&bull;"), []byte("&#8226;")},
		{[]byte("&nbsp;"), []byte(" ")},
		// Add more entity replacements as needed
	}

	for _, r := range replacements {
		data = bytes.ReplaceAll(data, r.old, r.new)
	}

	return data
}
