package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shivajichalise/rssagg/internal/database"
)

func scrape(
	db *database.Queries,
	concurrency int, // how many different go routines we want to do scraping on
	timeBetweenRequest time.Duration,
) {
	log.Printf("INFO: Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	/*
		for loop will execute immediately and then waits for the interval on the ticker
		if we had done:
		for range ticker.C {}
		then we'd actually wait for a minute upfront
	*/

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("ERROR: Could not fetch feeds: ", err)
			continue
		}

		waitGroup := &sync.WaitGroup{}

		for _, feed := range feeds {
			waitGroup.Add(1)

			go scrapeFeed(db, waitGroup, feed)
		}
		waitGroup.Wait()
	}
}

func scrapeFeed(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("ERROR: Could not mark feed as fetched: ", err)
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("ERROR: Could not fetch feed: ", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("ERROR: Could not parse published at date for item", item.Title, "on feed", feed.Name)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("ERROR: Failed to create post: ", err)
		}
	}
	log.Printf("INFO: Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
