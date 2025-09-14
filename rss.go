package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error construcing request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error getting resposne: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading response body: %w", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("error decoding xml data: %w", err)
	}
	unescapeFeed(&rssFeed)
	return &rssFeed, nil
}

func unescapeFeed(f *RSSFeed) {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	for i := range f.Channel.Item {
		f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
		f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
	}
}

func scrapeFeeds(db *database.Queries) error {
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	markedFeed, err := db.MarkFeedFetched(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error marking feed fetched: %w", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf(" - %s\n", item.Title)
	}

	return nil
}
