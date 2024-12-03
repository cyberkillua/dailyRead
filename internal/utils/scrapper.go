package utils

import (
	"context"
	"database/sql"
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
		return
	}

	log.Printf("Scrapped feed %v", page.Url)

	for _, item := range rss.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}
		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Error parsing pubDate: %v", err)
			continue
		}
		publishedAt := sql.NullTime{}
		if item.PubDate != "" {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		}
		post, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("Post %v already exists", post.Url)
				continue
			}
			log.Printf("Error creating post: %v", err)
			return
		}

		log.Printf("Created post %v", post.Url)
	}

}
