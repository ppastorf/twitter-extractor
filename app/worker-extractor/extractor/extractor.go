package extractor

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	batch "worker-extractor/extractor/tweet_batch"

	"github.com/dghubble/go-twitter/twitter"
)

type configSearchParams struct {
	Query           string `yaml:"query,omitempty"`
	Geocode         string `yaml:"geocode,omitempty"`
	Lang            string `yaml:"lang,omitempty"`
	Locale          string `yaml:"locale,omitempty"`
	ResultType      string `yaml:"result_type,omitempty"`
	Count           int    `yaml:"count,omitempty"`
	SinceID         int64  `yaml:"since_id,omitempty"`
	MaxID           int64  `yaml:"max_id,omitempty"`
	Until           string `yaml:"until,omitempty"`
	Since           string `yaml:"since,omitempty"`
	Filter          string `yaml:"filter,omitempty"`
	IncludeEntities *bool  `yaml:"include_entities,omitempty"`
	TweetMode       string `yaml:"tweet_mode,omitempty"`
}

func (confParams *configSearchParams) getAsApiParams() (*twitter.SearchTweetParams, error) {
	apiParams := &twitter.SearchTweetParams{
		Query:           confParams.Query,
		Geocode:         confParams.Geocode,
		Lang:            confParams.Lang,
		Locale:          confParams.Locale,
		ResultType:      confParams.ResultType,
		Count:           confParams.Count,
		SinceID:         confParams.SinceID,
		MaxID:           confParams.MaxID,
		Until:           confParams.Until,
		Since:           confParams.Since,
		Filter:          confParams.Filter,
		IncludeEntities: confParams.IncludeEntities,
		TweetMode:       confParams.TweetMode,
	}
	return apiParams, nil
}

type ExtractorConfig struct {
	SearchParams *configSearchParams `yaml:"searchParams"`
	Interval     string              `yaml:"interval"`
	LogTweets    bool                `yaml:"log_tweet_content"`
}

type (
	Extractor interface {
		GetTweetBatch()
		Start(chan<- batch.TweetBatch)
		End()
		processExtractions(chan<- batch.TweetBatch)
	}
	extractor struct {
		apiClient    *twitter.Client
		SearchParams *twitter.SearchTweetParams
		Interval     time.Duration
		LogTweets    bool
		end          chan struct{}
		extracted    chan []twitter.Tweet
		errs         chan error
	}
)

func NewExtractor(creds *TwitterCredentials, config *ExtractorConfig) (Extractor, error) {
	client, err := getApiClient(creds)
	if err != nil {
		log.Println("failed getting Twitter client: ", err)
		return nil, err
	}

	log_tweets := config.LogTweets

	searchParams, err := config.SearchParams.getAsApiParams()
	if err != nil {
		log.Println("invalid 'searchParams' from config: ", err)
		return nil, err
	}

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		log.Println("invalid 'interval' from config: ", err)
		return nil, err
	}

	return &extractor{
		apiClient:    client,
		SearchParams: searchParams,
		Interval:     interval,
		LogTweets:    log_tweets,
		extracted:    make(chan []twitter.Tweet),
		errs:         make(chan error),
	}, nil
}

func (e *extractor) GetTweetBatch() {
	log.Printf("Extractor: new job: GetTweetBatch\n")
	searchApiBatch, _, err := e.apiClient.Search.Tweets(e.SearchParams)
	if err != nil {
		log.Printf("Extractor: failed getting tweets: %v\n", err)
		e.errs <- err
		return
	}
	log.Printf("Extractor: job GetTweetBatch success\n")
	e.extracted <- searchApiBatch.Statuses
}

func (e *extractor) Start(dest chan<- batch.TweetBatch) {
	go e.processExtractions(dest)

	ticker := time.NewTicker(e.Interval)
	for {
		select {
		case <-ticker.C:
			go e.GetTweetBatch()
		case <-e.end:
			log.Printf("Extractor: terminating.")
			ticker.Stop()
			return
		}
	}
}

func (e *extractor) End() {
	close(e.end)
}

func genTweetBatch(apiReturn []twitter.Tweet) (*batch.TweetBatch, error) {
	itemsJson, err := json.Marshal(apiReturn)
	if err != nil {
		return nil, err
	}

	var items []batch.Tweet
	err = json.Unmarshal(itemsJson, &items)
	if err != nil {
		return nil, err
	}

	b := batch.TweetBatch{
		Id:          int64(1),
		ExtractedAt: "now",
		Size:        len(items),
		Items:       items,
	}
	return &b, nil
}

func (e *extractor) processExtractions(dest chan<- batch.TweetBatch) {
	for {
		select {
		case items := <-e.extracted:
			log.Printf("Extractor: extracted %d tweets.\n", len(items))
			if e.LogTweets {
				for _, t := range items {
					fmt.Printf("-> %s disse:\n%s\n\n", t.User.Name, t.FullText)
				}
			}
			newBatch, err := genTweetBatch(items)
			if err != nil {
				log.Printf("Extractor: failed generating tweet batch: %v\n", err)
			} else {
				log.Printf("Extractor: generated tweet batch.\n")
				dest <- *newBatch
			}
		case <-e.end:
			return
		}
	}
}
