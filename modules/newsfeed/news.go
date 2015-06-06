package newsfeed

import (
	"time"

	"github.com/jamesclonk-io/stdlib/web/newsreader"
)

func UpdateFeeds(news *newsreader.NewsReader) {
	news.InitializeFeeds()

	ticker := time.NewTicker(time.Hour * 1)
	go func(news *newsreader.NewsReader) {
		for {
			select {
			case <-ticker.C:
				news.UpdateFeeds()
			}
		}
	}(news)
}
