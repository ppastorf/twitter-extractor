package extractor

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterCredentials struct {
	ConsumerKey       string `yaml:"CONSUMER_KEY"`
	ConsumerSecret    string `yaml:"CONSUMER_KEY_SECRET"`
	AccessToken       string `yaml:"ACCESS_TOKEN"`
	AccessTokenSecret string `yaml:"ACCESS_TOKEN_SECRET"`
}

func getApiClient(creds *TwitterCredentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}
