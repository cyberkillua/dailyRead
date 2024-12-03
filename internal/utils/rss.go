package utils

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSS struct {
	Title       string    `xml:"channel>title"`
	Link        string    `xml:"channel>link"`
	Description string    `xml:"description"`
	Language    string    `xml:"channel>language"`
	Items       []RSSItem `xml:"channel>item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
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

	rssFeed := RSS{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSS{}, err
	}

	return rssFeed, nil
}
