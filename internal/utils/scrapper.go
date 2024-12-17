package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/google/uuid"
)

func StartScrapping(db *database.Queries, concurrency int, durationBetween time.Duration) {
	log.Printf("Scarping on %v goroutines every %v", concurrency, durationBetween)

	ticker := time.NewTicker(durationBetween)
	for ; ; <-ticker.C {
		pages, err := db.GetNextWebpageToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Printf("Error getting feeds to scrap: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}

		for _, page := range pages {
			wg.Add(1)

			go scrapeFeed(db, wg, page)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, page database.Webpage) {
	defer wg.Done()

	_, err := db.MarkWebpageAsFetched(context.Background(), page.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rss, err := urlToRSS(page.Url)
	if err != nil {
		log.Printf("Error scrapping feed: %v", err)
		log.Printf("This is the problematic: %v", page.Url)
		return
	}

	log.Printf("Scrapped feed %v", page.Url)
	log.Printf("Found %v channels", len(rss.Channel.Items))

	// log.Printf("All Rss Items: %v", rss)

	for _, item := range rss.Channel.Items {

		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		pageName := sql.NullString{}
		if page.Name != "" {
			pageName = sql.NullString{String: item.Title, Valid: true}
		}

		publishedAt := sql.NullTime{}
		if item.PubDate != "" {
			t, err := parseDate(item.PubDate)
			if err != nil {
				log.Printf("Error parsing pubDate: %v", err)
				log.Printf("Undefined pubDate for item %v", item.PubDate)
				continue
			}

			// Successfully parsed pubDate
			publishedAt = sql.NullTime{Time: t, Valid: true}
		}

		if item.Link == "" {
			log.Printf("Undefined link for item %v", item.Title)
			continue
		}
		post, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: publishedAt,
			Postname:    pageName,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				// log.Printf("Post %v already exists", post.Url)
				continue
			}
			log.Printf("Error creating post: %v", err)
			return
		}

		log.Printf("Created post %v", post.Url)
	}

}

func parseDate(pubDate string) (time.Time, error) {
	// Define the possible date formats
	formats := []string{
		time.RFC1123,  // Example: "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z, // Example: "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,  // Example: "2006-01-02T15:04:05Z07:00"
		time.RFC822,   // Example: "02 Jan 06 15:04 MST"
		time.RFC822Z,  // Example: "02 Jan 06 15:04 -0700"
		time.RFC850,   // Example: "Monday, 02-Jan-06 15:04:05 MST"
		time.RubyDate, // Example: "Mon Jan 02 15:04:05 -0700 2006"
	}

	// Iterate through the formats and try parsing
	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return t, nil // Successfully parsed
		}
	}

	// If all formats fail, return an error
	return time.Time{}, fmt.Errorf("unknown date format: %s", pubDate)
}
