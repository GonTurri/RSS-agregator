package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/GonTurri/RSS-agregator/internal/database"
)

// representa un xml
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetch
func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return RSSFeed{}, err
	}

	result := RSSFeed{}

	err = xml.Unmarshal(data, &result)

	if err != nil {
		return RSSFeed{}, err
	}

	return result, nil
}

func startScraping(db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(),
			int32(concurrency))

		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(wg, db, feed)
		}

		wg.Wait()

	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFetchedFeed(context.Background(), feed.ID)

	if err != nil {
		log.Println("error marking feed as fetched: ", err)
		return
	}

	RssFeed, err := urlToFeed(feed.Url)

	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range RssFeed.Channel.Item {
		log.Println("found post: ", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(RssFeed.Channel.Item))

}
