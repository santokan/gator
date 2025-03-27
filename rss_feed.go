package main

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed fetches an RSS feed from the given URL and returns a filled-out RSSFeed struct.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Perform the HTTP request
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch feed: " + resp.Status)
	}

	// Decode the response body into the RSSFeed struct
	var feed RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}

	// Unescape HTML entities in the Title and Description fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
