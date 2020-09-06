package newsfeed

import (
	"fmt"
	"time"

	"github.com/jamesclonk-io/stdlib/web/newsreader"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	newsfeedUpdates = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jcio_frontend_newsfeed_updates",
		Help: "Total number of JCIO frontend newsfeed updates.",
	})
	newsfeedFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jcio_frontend_newsfeed_failures",
		Help: "Total number of JCIO frontend newsfeed failures/panics.",
	})
)

func UpdateFeeds(news *newsreader.NewsReader) {
	news.InitializeFeeds()

	ticker := time.NewTicker(time.Hour * 1)
	go func(news *newsreader.NewsReader) {
		for {
			select {
			case <-ticker.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("recovering from panic in UpdateFeeds/news.UpdateFeeds(), error is: %v\n", r)
							newsfeedFailures.Inc()
						}
					}()
					news.UpdateFeeds()
					newsfeedUpdates.Inc()
				}()
			}
		}
	}(news)
}
