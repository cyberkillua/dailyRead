package utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// type RSS struct {
// 	Title       string    `xml:"channel>title"`
// 	Link        string    `xml:"channel>link"`
// 	Description string    `xml:"description,omitempty"`
// 	Language    string    `xml:"channel>language,omitempty"`
// 	Items       []RSSItem `xml:"channel>item"`
// }

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
	Description string `xml:"summary"` // Use "summary" for description
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
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSS{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSS{}, err
	}
	// Detect feed type by unmarshalling into a generic structure
	var root struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(data, &root); err != nil {
		return RSS{}, fmt.Errorf("failed to parse XML: %w", err)
	}

	switch strings.ToLower(root.XMLName.Local) {
	case "rss":
		// Parse RSS feed
		var rssFeed RSS
		if err := xml.Unmarshal(data, &rssFeed); err != nil {
			return RSS{}, fmt.Errorf("failed to parse RSS feed: %w", err)
		}
		return rssFeed, nil
	case "feed":
		// Parse Atom feed
		var atomFeed Atom
		if err := xml.Unmarshal(data, &atomFeed); err != nil {
			return RSS{}, fmt.Errorf("failed to parse Atom feed: %w", err)
		}
		// Convert Atom to RSS-like structure
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
