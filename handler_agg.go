package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/santokan/gator/internal/database"
)

func handlerAgg(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: agg <interval>")
	}

	timeBetweenRequests, err := time.ParseDuration(c.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing interval: %w", err)
	}

	log.Printf("Collecting feeds every %v...", timeBetweenRequests)

	// create a ticker
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s.db)
		if err != nil {
			log.Printf("Error scraping feeds: %v", err)
		}
	}
}

func scrapeFeeds(s *database.Queries) error {
	// get the next feed from DB
	feed, err := s.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed: %w", err)
	}

	fmt.Println("Found feed to fetch: ", feed.Name)

	// mark it as fetched
	err = s.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("unable to mark the feed '%s' as fetched: %w", feed.Name, err)
	}

	// fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed '%s' with url '%s': %w", feed.Name, feed.Url, err)
	}

	// iterate over the feeds
	for _, item := range rssFeed.Channel.Item {
		fmt.Println("Title:", item.Title)
	}

	fmt.Printf("Feed '%s' collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

	return nil
}
