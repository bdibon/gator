package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't build request: %w", err)
	}
	req.Header.Add("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	feed := &RSSFeed{}
	err = xml.Unmarshal(body, feed)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	decodeHTMLEntities(feed)
	return feed, nil
}

func decodeHTMLEntities(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
}
